# AML Screening AI Workflow — Discussion Summary

## Context

Reviewed the `world_check_process` function in `world_check_router.py`. The current workflow bulk-sends all screening hits with case data to an internal LLM service (`worldcheck_compliance_analyst` prompt via LiteLLM) and receives back an array of AI recommendations. Processing happens in batches of 25 hits.

## Problems with Current Bulk Approach

1. **Lost in the Middle** — LLM quality degrades on items in the middle of long lists. Hit #13 in a batch of 25 gets less attention than #1 or #25.
2. **No per-result reasoning trace** — no audit trail of why the LLM made a specific recommendation. Compliance officers can't verify decisions.
3. **All-or-nothing failure** — one malformed hit or LLM error fails the entire batch. No partial recovery.
4. **No evaluation loop** — LLM output goes straight to the response with no validation.

## Recommended Techniques

| Technique | Purpose |
|---|---|
| Map-Reduce Prompting | Analyze each hit individually (map), then aggregate (reduce) |
| Chain-of-Thought (CoT) | Force LLM to show step-by-step reasoning per hit |
| Structured Output / JSON Schema | Guarantee output shape, reject malformed responses |
| LLM-as-Judge / Evals | Second pass to validate recommendation quality |
| Guardrails (NeMo, Guardrails AI) | Input/output validation around LLM calls |
| Confidence Scoring | LLM assigns confidence per recommendation; low confidence routes to human review |
| Observability (LangSmith, Langfuse) | Trace prompts, responses, latency, token usage |
| Semantic Caching | Reuse past recommendations for similar entities/hits |
| Human-in-the-Loop (HITL) | Flag uncertain recommendations for manual review |

## Proposed Architecture: Fan-out / Fan-in

```
Event from RabbitMQ
  -> Data Incoming (case + hits[])
    -> Worker 1 (hit[0]) --\
    -> Worker 2 (hit[1]) ---\
    -> Worker 3 (hit[2]) -----> Update Parent Data to DB -> Done
    -> Worker N (hit[N]) ---/
```

Each worker:
1. Data processing (mapping, cleaning)
2. AI recommendation (per-hit CoT prompt)
3. Insert per-hit result to DB

## Go errgroup vs RabbitMQ per Hit

**Decision: Use Go errgroup for fan-out within a single screening request.**

The bottleneck is the LLM call (1-5s per hit, I/O-bound). Go goroutines waiting on HTTP responses cost almost nothing. RabbitMQ is already the job-level transport — fanning out each hit to RabbitMQ is over-engineering unless you need multi-pod processing or hit-level retry that survives pod restarts.

```go
func processHits(ctx context.Context, caseRecord CaseRecord, hits []Hit) ([]Recommendation, error) {
    results := make([]Recommendation, len(hits))
    g, ctx := errgroup.WithContext(ctx)
    sem := make(chan struct{}, 10) // bound concurrency to LLM rate limit

    for i, hit := range hits {
        g.Go(func() error {
            sem <- struct{}{}
            defer func() { <-sem }()
            rec, err := callLLM(ctx, caseRecord, hit)
            if err != nil {
                return fmt.Errorf("hit %d: %w", i, err)
            }
            results[i] = rec
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return results, err
    }
    return results, nil
}
```

## Resource Concern: 200 Hits Per Case

### Rate Limiting
Bound concurrency with semaphore matched to LiteLLM rate limit. 200 hits with `sem=10` and 20 RPM provider limit = ~20-40 seconds total.

### Cost: ~10x More Tokens
Per-hit repeats case context in every call. Mitigation: **hybrid pre-filter**.

```go
// Rule-based pre-filter before LLM
if hit.MatchStrength == "EXACT" && hit.Category == "SANCTIONS" {
    // obvious true match — skip LLM
}
if hit.MatchScore < 0.3 {
    // obvious false positive — skip LLM
}
// only ambiguous hits go to LLM
```

In practice, ~30-50 out of 200 hits typically need LLM analysis.

## CoT Prompt Structure (per hit)

The prompt forces 4 sequential analysis steps:

1. **Name Analysis** — compare subject name vs hit aliases (exact, transliteration, pinyin/hanzi, telegraph codes). Output: EXACT / STRONG / PARTIAL / WEAK.
2. **Biographic Comparison** — compare DOB, gender, nationality, location. Each field: MATCH / CONFLICT / UNAVAILABLE. Explicitly flags missing data as a gap, not confirmation.
3. **Risk Assessment** — nature of offense, recency, sentencing, associate network size, regulatory implications.
4. **Conclusion** — verdict based on steps 1-3.

### Output Schema

```json
{
  "reasoning": {
    "name_analysis": "...",
    "name_match_strength": "EXACT | STRONG | PARTIAL | WEAK",
    "biographic_comparison": "...",
    "biographic_fields": {
      "dob": "MATCH | CONFLICT | UNAVAILABLE",
      "gender": "MATCH | CONFLICT | UNAVAILABLE",
      "nationality": "MATCH | CONFLICT | UNAVAILABLE",
      "location": "MATCH | CONFLICT | UNAVAILABLE"
    },
    "risk_assessment": "..."
  },
  "verdict": "TRUE_MATCH | FALSE_POSITIVE | NEEDS_REVIEW",
  "confidence": 0.0,
  "recommended_action": "...",
  "key_factors": ["..."]
}
```

### Prompt Design Notes
- Telegraph code aliases (e.g. `3769 0193 2494` for `王俊明`) surfaced explicitly to prevent LLM misinterpretation
- Comparison table forces per-field verdicts — prevents lazy "looks similar enough" reasoning
- Missing DOB caveat called out explicitly — LLMs tend to treat missing data as "not conflicting"
- `recommended_action` gives compliance officers actionable next steps
- `key_factors` array enables UI display of decision drivers without parsing free-text
