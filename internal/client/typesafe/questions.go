package typesafe

// ReasoningQuestions contains the 10 reasoning questions for the TypeSafe API
var ReasoningQuestions = map[string]any{
	"name_match_strength": map[string]any{
		"type": "score",
		"instructions": "Compare the subject name against all hit names and aliases. Consider: exact character match, pinyin-hanzi transliteration, romanization variants (Wade-Giles, Yale, etc.), partial matches, and telegraph code correspondence. Score the strongest match found across all aliases. When torn between two levels, pick the lower.",
		"criteria": []string{
			"No meaningful similarity between subject name and any hit name or alias. Different characters, different phonetics, no plausible transliteration link.",
			"Surface-level overlap only: shared surname OR shared given name but not both, a common character appearing by coincidence, or a partial phonetic similarity that could apply to many names.",
			"Clear phonetic or transliteration link between subject and hit names but not exact: consistent pinyin-hanzi mapping with minor romanization variant (e.g. Cheung vs Zhang), name order transposition, or one matching alias among several.",
			"Exact character-for-character match in at least one alias, confirmed by consistent transliteration across all name forms (pinyin, telegraph code if present). No ambiguity in the name link.",
		},
	},
	"name_distinguishing_power": map[string]any{
		"type": "choice",
		"instructions": "How common is this name in the subject's cultural and national context? A common name matching exactly carries less weight than a rare or distinctive name matching.",
		"criteria": map[string]string{
			"high_frequency":       "Extremely common name in the cultural context (e.g. Zhang Wei, Li Ming, John Smith). Exact match alone is insufficient to distinguish individuals.",
			"moderate_frequency":   "Moderately common name. Match is meaningful but not uniquely identifying without biographic corroboration.",
			"low_frequency":        "Uncommon or distinctive name. An exact match is strong evidence of identity on its own.",
			"compound_distinctive": "Multi-part name, unusual transliteration, or name combination that is highly unlikely to be coincidental.",
		},
	},
	"dob_comparison": map[string]any{
		"type": "choice",
		"instructions": "Compare the subject's date of birth against the hit's date of birth from the comparison fields. A missing DOB on either side is a data gap — score it as unavailable, never as confirmation or denial.",
		"criteria": map[string]string{
			"exact_match":   "Full date (YYYY-MM-DD) matches exactly.",
			"partial_match": "Year matches but month/day differs or is missing. Or month/year match but day is absent.",
			"conflict":      "Both dates are present and differ by more than a minor transcription error (e.g. different year).",
			"unavailable":   "DOB is missing or null on one or both sides. Cannot compare.",
		},
	},
	"gender_comparison": map[string]any{
		"type": "choice",
		"instructions": "Compare the subject's gender against the hit's gender.",
		"criteria": map[string]string{
			"match":       "Both genders are present and identical.",
			"conflict":    "Both genders are present and different.",
			"unavailable": "Gender is missing on one or both sides.",
		},
	},
	"nationality_comparison": map[string]any{
		"type": "choice",
		"instructions": "Compare the subject's nationality/citizenship against the hit's nationality. Consider country code mappings (e.g. CHN = China).",
		"criteria": map[string]string{
			"match":       "Nationalities refer to the same country.",
			"conflict":    "Nationalities refer to different countries.",
			"unavailable": "Nationality is missing on one or both sides.",
		},
	},
	"location_comparison": map[string]any{
		"type": "choice",
		"instructions": "Compare the subject's country/location/place of birth against the hit's addresses and location details. Country-level match counts even if city-level detail is missing.",
		"criteria": map[string]string{
			"match":       "Locations refer to the same country or city.",
			"conflict":    "Locations refer to different countries.",
			"partial":     "Same country but different or missing city, or only one side has location data at country level.",
			"unavailable": "Location is missing on one or both sides.",
		},
	},
	"biographic_conflict_count": map[string]any{
		"type": "score",
		"instructions": "Count how many biographic fields (DOB, gender, nationality, location) have a definitive CONFLICT — meaning both sides have data and it does not match. Unavailable fields do not count as conflicts.",
		"criteria": []string{
			"No biographic conflicts. All available fields match or are unavailable.",
			"Exactly 1 biographic field conflicts.",
			"Exactly 2 biographic fields conflict.",
			"3 or more biographic fields conflict.",
		},
	},
	"risk_category_severity": map[string]any{
		"type": "score",
		"instructions": "Rate the regulatory severity of the hit's category and source. Consider the nature of the listing and its AML/CFT implications.",
		"criteria": []string{
			"Low-severity category with minimal AML/CFT implications: media mention, historical PEP with no adverse info, expired minor regulatory action.",
			"Moderate severity: active PEP, post-conviction for non-financial crime, or inclusion in a regional watchlist without sanctions.",
			"High severity: sanctions list (OFAC, EU, UN), conviction for financial crime (fraud, money laundering, terrorist financing), or active law enforcement investigation.",
			"Critical severity: primary sanctions designation (SDN, OFSI) for terrorism financing or proliferation, active Interpol notice, or designated by multiple jurisdictions.",
		},
	},
	"adverse_info_recency": map[string]any{
		"type": "choice",
		"instructions": "How recent is the most recent adverse information or listing action? Judge from dates mentioned in the adverse info, report text, and source update timestamps.",
		"criteria": map[string]string{
			"current":  "Within the last 2 years. Active and likely still in force.",
			"recent":   "2-5 years ago. May still be relevant but could have changed.",
			"dated":    "More than 5 years ago. Circumstances may have materially changed.",
			"no_date":  "No dates provided in adverse information. Cannot assess recency.",
		},
	},
	"associate_network_risk": map[string]any{
		"type": "choice",
		"instructions": "Assess the risk signal from the hit's associate/connection network. A large network of sanctioned or convicted associates amplifies risk.",
		"criteria": map[string]string{
			"no_associates": "No associates or connections listed.",
			"low_risk":      "Associates listed but none are sanctioned, convicted, or otherwise flagged.",
			"moderate_risk": "Some associates are themselves flagged, PEP-connected, or in related industries.",
			"high_risk":     "Multiple associates are sanctioned, convicted, or part of a known criminal/terrorist network.",
		},
	},
}

// ClassificationQuestions contains the 4 classification questions for the TypeSafe API
var ClassificationQuestions = map[string]any{
	"verdict": map[string]any{
		"type": "choice",
		"instructions": "Based on the totality of name match strength, biographic field comparisons, risk severity, and available evidence: classify this screening hit. A strong name match with missing biographic data is NOT a confirmed match — it requires review. Multiple biographic conflicts override even an exact name match.",
		"criteria": map[string]string{
			"TRUE_MATCH":    "Name match is STRONG or EXACT, AND no biographic conflicts exist, AND at least 2 biographic fields positively match. Evidence collectively confirms this is the same individual.",
			"FALSE_POSITIVE": "One or more definitive biographic conflicts (different DOB, different gender, different nationality) that cannot be explained by data quality issues. OR name match is WEAK with no supporting biographic evidence.",
			"NEEDS_REVIEW":  "Name match is STRONG or EXACT but biographic data is insufficient to confirm or deny: key fields (especially DOB) are missing/unavailable, or there is a mix of matching and ambiguous data. Human analyst must obtain additional information.",
		},
	},
	"confidence_level": map[string]any{
		"type": "score",
		"instructions": "How confident is the verdict? High confidence means the available evidence clearly supports the conclusion with little ambiguity. Low confidence means the data is thin, contradictory, or borderline.",
		"criteria": []string{
			"Very low confidence. Data is sparse or contradictory. The verdict is a best guess.",
			"Low confidence. Some supporting evidence but significant gaps or ambiguity remain.",
			"Moderate confidence. The verdict is supported by evidence but one or two factors introduce doubt.",
			"High confidence. Multiple independent data points align to support the verdict with minimal ambiguity.",
			"Very high confidence. All available evidence consistently and unambiguously supports the verdict.",
		},
	},
	"escalation_required": map[string]any{
		"type": "noul",
		"instructions": "Does this hit require escalation to a senior compliance officer or MLRO beyond standard resolution?",
		"criteria": map[string]string{
			"true":  "The hit involves a sanctions designation, terrorism financing, or proliferation — or the verdict is TRUE_MATCH regardless of category.",
			"false": "The hit can be resolved at analyst level: clear false positive, or a NEEDS_REVIEW that only requires additional data gathering.",
		},
	},
	"additional_info_needed": map[string]any{
		"type": "noul",
		"instructions": "Is additional identifying information required to reach a definitive verdict?",
		"criteria": map[string]string{
			"true":  "One or more critical biographic fields are unavailable and the verdict cannot be confirmed without them.",
			"false": "Available data is sufficient to reach a definitive TRUE_MATCH or FALSE_POSITIVE conclusion.",
		},
	},
}
