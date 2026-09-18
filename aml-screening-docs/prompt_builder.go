package screening

import (
	"fmt"
	"strings"
)

// Subject represents the customer under review from your KYC/CRM system.
type Subject struct {
	Name        string // Full name, native script preferred (e.g. 王俊明)
	DOB         string // YYYY-MM-DD
	Nationality string // ISO 3166-1 alpha-3 (e.g. CHN)
	Gender      string // MALE | FEMALE
	EntityType  string // INDIVIDUAL | ENTITY
}

// Hit represents a single World Check screening result.
type Hit struct {
	ReferenceID  string            // e.g. e_tr_wci_4995308
	PrimaryName  string            // Main name on watchlist record
	MatchedTerm  string            // Alias/variant that triggered the match
	Provider     string            // WATCHLIST | PEP | SANCTION
	Category     string            // e.g. POST CONVICTION, PEP, SANCTION LIST
	Gender       string
	Nationality  string
	Location     string
	Aliases      []Alias
	Comparisons  []ComparisonField
	AdverseInfo  AdverseInfo
	Associates   []Associate
	Sources      []string
}

type Alias struct {
	Name     string // e.g. "Junming WANG"
	Type     string // PRIMARY | AKA | NATIVE_AKA | DBA | FORMERLY_KNOWN_AS
	Language string // optional, e.g. "Chinese"
}

type ComparisonField struct {
	Field     string // e.g. NATIVE_AKA, DATE_OF_BIRTH
	Submitted string
	Matched   string
	Result    string // EXACT | CLOSE | NOT_MATCHED | N/A
}

type AdverseInfo struct {
	Categories []string
	Keywords   []string
	Reports    []string
}

type Associate struct {
	Name       string
	Category   string
	EntityType string
}

const systemPrompt = `You are an AML/CFT compliance analyst. Analyze one World Check screening hit against the subject and determine: TRUE_MATCH, FALSE_POSITIVE, or NEEDS_REVIEW.`

// SystemPrompt returns the static system prompt.
func SystemPrompt() string {
	return systemPrompt
}

// BuildUserPrompt assembles the user prompt from subject and hit data.
func BuildUserPrompt(s Subject, h Hit) string {
	var b strings.Builder

	// Subject block
	fmt.Fprintf(&b, "SUBJECT\nName: %s\nDOB: %s\nNationality: %s\nGender: %s\nType: %s\n",
		s.Name, s.DOB, s.Nationality, s.Gender, s.EntityType)

	// Hit block
	fmt.Fprintf(&b, "\nHIT [%s]\nName: %s\nMatched: %s\nProvider: %s\nCategory: %s\nGender: %s\nNationality: %s\nLocation: %s\n",
		h.ReferenceID, h.PrimaryName, h.MatchedTerm, h.Provider, h.Category,
		fieldOrNA(h.Gender), fieldOrNA(h.Nationality), fieldOrNA(h.Location))

	// Aliases
	if len(h.Aliases) > 0 {
		b.WriteString("\nAliases:\n")
		for i, a := range h.Aliases {
			if a.Language != "" {
				fmt.Fprintf(&b, "%d. %s (%s, %s)\n", i+1, a.Name, a.Type, a.Language)
			} else {
				fmt.Fprintf(&b, "%d. %s (%s)\n", i+1, a.Name, a.Type)
			}
		}
	}

	// Comparison fields
	if len(h.Comparisons) > 0 {
		b.WriteString("\nField Comparison:\n")
		for _, c := range h.Comparisons {
			fmt.Fprintf(&b, "- %s: submitted=%s matched=%s result=%s\n",
				c.Field, fieldOrNA(c.Submitted), fieldOrNA(c.Matched), c.Result)
		}
	}

	// Adverse information
	ai := h.AdverseInfo
	if len(ai.Categories) > 0 || len(ai.Keywords) > 0 || len(ai.Reports) > 0 {
		b.WriteString("\nAdverse Info:\n")
		if len(ai.Categories) > 0 {
			fmt.Fprintf(&b, "- Categories: %s\n", strings.Join(ai.Categories, ", "))
		}
		if len(ai.Keywords) > 0 {
			fmt.Fprintf(&b, "- Keywords: %s\n", strings.Join(ai.Keywords, ", "))
		}
		for _, r := range ai.Reports {
			fmt.Fprintf(&b, "- Report: %s\n", r)
		}
	}

	// Associates
	if len(h.Associates) > 0 {
		fmt.Fprintf(&b, "\nAssociates (%d):\n", len(h.Associates))

		// Group by category
		grouped := make(map[string][]string)
		for _, a := range h.Associates {
			grouped[a.Category] = append(grouped[a.Category], a.Name)
		}
		for cat, names := range grouped {
			fmt.Fprintf(&b, "- %s: %s\n", cat, strings.Join(names, ", "))
		}
	}

	// Sources
	if len(h.Sources) > 0 {
		b.WriteString("\nSources:\n")
		for _, src := range h.Sources {
			fmt.Fprintf(&b, "- %s\n", src)
		}
	}

	// Instructions
	b.WriteString(`
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
}`)

	return b.String()
}

func fieldOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
