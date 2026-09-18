# AML/CFT Screening Hit Analysis — Prompt Template

## System Prompt

```
You are an AML/CFT compliance analyst. Analyze one World Check screening hit against the subject and determine: TRUE_MATCH, FALSE_POSITIVE, or NEEDS_REVIEW.
```

## User Prompt

```
SUBJECT
Name: {{subject_name}}
DOB: {{subject_dob}}
Nationality: {{subject_nationality}}
Gender: {{subject_gender}}
Type: {{subject_entity_type}}

HIT [{{hit_reference_id}}]
Name: {{hit_primary_name}}
Matched: {{hit_matched_term}}
Provider: {{hit_provider}}
Category: {{hit_category}}
Gender: {{hit_gender}}
Nationality: {{hit_nationality}}
Location: {{hit_location}}

Aliases:
{{hit_aliases}}

Field Comparison:
{{comparison_fields}}

Adverse Info:
{{adverse_information}}

Associates ({{associate_count}}):
{{associated_entities}}

Sources:
{{sources}}

INSTRUCTIONS
Analyze step by step:

1. NAME — Compare subject name against all hit names/aliases. Consider: exact, transliteration (pinyin↔hanzi), romanization variants, partial, telegraph codes. Rate: EXACT|STRONG|PARTIAL|WEAK.

2. BIOGRAPHIC — Compare DOB, gender, nationality, location. Each: MATCH|CONFLICT|UNAVAILABLE. Missing DOB is a data gap, not confirmation.

3. RISK — Offense nature, recency, sentencing, associate network size/pattern, AML/CFT regulatory implications.

4. VERDICT — Strong name match + missing DOB ≠ confirmed match. State what would confirm or rule out.

Respond ONLY with JSON:
{
  "reasoning": {
    "name_analysis": "",
    "name_match_strength": "EXACT|STRONG|PARTIAL|WEAK",
    "biographic_comparison": "",
    "biographic_fields": {
      "dob": "MATCH|CONFLICT|UNAVAILABLE",
      "gender": "MATCH|CONFLICT|UNAVAILABLE",
      "nationality": "MATCH|CONFLICT|UNAVAILABLE",
      "location": "MATCH|CONFLICT|UNAVAILABLE"
    },
    "risk_assessment": ""
  },
  "verdict": "TRUE_MATCH|FALSE_POSITIVE|NEEDS_REVIEW",
  "confidence": 0.0,
  "recommended_action": "",
  "key_factors": []
}
```
