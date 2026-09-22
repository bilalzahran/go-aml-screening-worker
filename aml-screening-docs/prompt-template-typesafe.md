# AML/CFT Screening Hit Analysis — TypeSafe Jev Prompt Template

## Context

```json
{
  "subject": {
    "name": "{{subject_name}}",
    "dob": "{{subject_dob}}",
    "nationality": "{{subject_nationality}}",
    "gender": "{{subject_gender}}",
    "entity_type": "{{subject_entity_type}}"
  },
  "hit": {
    "reference_id": "{{hit_reference_id}}",
    "primary_name": "{{hit_primary_name}}",
    "matched_term": "{{hit_matched_term}}",
    "provider": "{{hit_provider}}",
    "category": "{{hit_category}}",
    "gender": "{{hit_gender}}",
    "nationality": "{{hit_nationality}}",
    "location": "{{hit_location}}",
    "aliases": {{hit_aliases}},
    "comparison_fields": {{comparison_fields}},
    "adverse_information": "{{adverse_information}}",
    "associate_count": {{associate_count}},
    "associated_entities": {{associated_entities}},
    "sources": {{sources}}
  }
}
```

## Questions — Reasoning Layer

```json
{
  "name_match_strength": {
    "type": "score",
    "instructions": "Compare the subject name against all hit names and aliases. Consider: exact character match, pinyin-hanzi transliteration, romanization variants (Wade-Giles, Yale, etc.), partial matches, and telegraph code correspondence. Score the strongest match found across all aliases. When torn between two levels, pick the lower.",
    "criteria": [
      "No meaningful similarity between subject name and any hit name or alias. Different characters, different phonetics, no plausible transliteration link.",
      "Surface-level overlap only: shared surname OR shared given name but not both, a common character appearing by coincidence, or a partial phonetic similarity that could apply to many names.",
      "Clear phonetic or transliteration link between subject and hit names but not exact: consistent pinyin-hanzi mapping with minor romanization variant (e.g. Cheung vs Zhang), name order transposition, or one matching alias among several.",
      "Exact character-for-character match in at least one alias, confirmed by consistent transliteration across all name forms (pinyin, telegraph code if present). No ambiguity in the name link."
    ]
  },
  "name_distinguishing_power": {
    "type": "choice",
    "instructions": "How common is this name in the subject's cultural and national context? A common name matching exactly carries less weight than a rare or distinctive name matching.",
    "criteria": {
      "high_frequency": "Extremely common name in the cultural context (e.g. Zhang Wei, Li Ming, John Smith). Exact match alone is insufficient to distinguish individuals.",
      "moderate_frequency": "Moderately common name. Match is meaningful but not uniquely identifying without biographic corroboration.",
      "low_frequency": "Uncommon or distinctive name. An exact match is strong evidence of identity on its own.",
      "compound_distinctive": "Multi-part name, unusual transliteration, or name combination that is highly unlikely to be coincidental."
    }
  },
  "dob_comparison": {
    "type": "choice",
    "instructions": "Compare the subject's date of birth against the hit's date of birth from the comparison fields. A missing DOB on either side is a data gap — score it as unavailable, never as confirmation or denial.",
    "criteria": {
      "exact_match": "Full date (YYYY-MM-DD) matches exactly.",
      "partial_match": "Year matches but month/day differs or is missing. Or month/year match but day is absent.",
      "conflict": "Both dates are present and differ by more than a minor transcription error (e.g. different year).",
      "unavailable": "DOB is missing or null on one or both sides. Cannot compare."
    }
  },
  "gender_comparison": {
    "type": "choice",
    "instructions": "Compare the subject's gender against the hit's gender.",
    "criteria": {
      "match": "Both genders are present and identical.",
      "conflict": "Both genders are present and different.",
      "unavailable": "Gender is missing on one or both sides."
    }
  },
  "nationality_comparison": {
    "type": "choice",
    "instructions": "Compare the subject's nationality/citizenship against the hit's nationality. Consider country code mappings (e.g. CHN = China).",
    "criteria": {
      "match": "Nationalities refer to the same country.",
      "conflict": "Nationalities refer to different countries.",
      "unavailable": "Nationality is missing on one or both sides."
    }
  },
  "location_comparison": {
    "type": "choice",
    "instructions": "Compare the subject's country/location/place of birth against the hit's addresses and location details. Country-level match counts even if city-level detail is missing.",
    "criteria": {
      "match": "Locations refer to the same country or city.",
      "conflict": "Locations refer to different countries.",
      "partial": "Same country but different or missing city, or only one side has location data at country level.",
      "unavailable": "Location is missing on one or both sides."
    }
  },
  "biographic_conflict_count": {
    "type": "score",
    "instructions": "Count how many biographic fields (DOB, gender, nationality, location) have a definitive CONFLICT — meaning both sides have data and it does not match. Unavailable fields do not count as conflicts.",
    "criteria": [
      "No biographic conflicts. All available fields match or are unavailable.",
      "Exactly 1 biographic field conflicts.",
      "Exactly 2 biographic fields conflict.",
      "3 or more biographic fields conflict."
    ]
  },
  "risk_category_severity": {
    "type": "score",
    "instructions": "Rate the regulatory severity of the hit's category and source. Consider the nature of the listing and its AML/CFT implications.",
    "criteria": [
      "Low-severity category with minimal AML/CFT implications: media mention, historical PEP with no adverse info, expired minor regulatory action.",
      "Moderate severity: active PEP, post-conviction for non-financial crime, or inclusion in a regional watchlist without sanctions.",
      "High severity: sanctions list (OFAC, EU, UN), conviction for financial crime (fraud, money laundering, terrorist financing), or active law enforcement investigation.",
      "Critical severity: primary sanctions designation (SDN, OFSI) for terrorism financing or proliferation, active Interpol notice, or designated by multiple jurisdictions."
    ]
  },
  "adverse_info_recency": {
    "type": "choice",
    "instructions": "How recent is the most recent adverse information or listing action? Judge from dates mentioned in the adverse info, report text, and source update timestamps.",
    "criteria": {
      "current": "Within the last 2 years. Active and likely still in force.",
      "recent": "2-5 years ago. May still be relevant but could have changed.",
      "dated": "More than 5 years ago. Circumstances may have materially changed.",
      "no_date": "No dates provided in adverse information. Cannot assess recency."
    }
  },
  "associate_network_risk": {
    "type": "choice",
    "instructions": "Assess the risk signal from the hit's associate/connection network. A large network of sanctioned or convicted associates amplifies risk.",
    "criteria": {
      "no_associates": "No associates or connections listed.",
      "low_risk": "Associates listed but none are sanctioned, convicted, or otherwise flagged.",
      "moderate_risk": "Some associates are themselves flagged, PEP-connected, or in related industries.",
      "high_risk": "Multiple associates are sanctioned, convicted, or part of a known criminal/terrorist network."
    }
  }
}
```

## Questions — Classification Layer

```json
{
  "verdict": {
    "type": "choice",
    "instructions": "Based on the totality of name match strength, biographic field comparisons, risk severity, and available evidence: classify this screening hit. A strong name match with missing biographic data is NOT a confirmed match — it requires review. Multiple biographic conflicts override even an exact name match.",
    "criteria": {
      "TRUE_MATCH": "Name match is STRONG or EXACT, AND no biographic conflicts exist, AND at least 2 biographic fields positively match. Evidence collectively confirms this is the same individual.",
      "FALSE_POSITIVE": "One or more definitive biographic conflicts (different DOB, different gender, different nationality) that cannot be explained by data quality issues. OR name match is WEAK with no supporting biographic evidence.",
      "NEEDS_REVIEW": "Name match is STRONG or EXACT but biographic data is insufficient to confirm or deny: key fields (especially DOB) are missing/unavailable, or there is a mix of matching and ambiguous data. Human analyst must obtain additional information."
    }
  },
  "confidence_level": {
    "type": "score",
    "instructions": "How confident is the verdict? High confidence means the available evidence clearly supports the conclusion with little ambiguity. Low confidence means the data is thin, contradictory, or borderline.",
    "criteria": [
      "Very low confidence. Data is sparse or contradictory. The verdict is a best guess.",
      "Low confidence. Some supporting evidence but significant gaps or ambiguity remain.",
      "Moderate confidence. The verdict is supported by evidence but one or two factors introduce doubt.",
      "High confidence. Multiple independent data points align to support the verdict with minimal ambiguity.",
      "Very high confidence. All available evidence consistently and unambiguously supports the verdict."
    ]
  },
  "escalation_required": {
    "type": "noul",
    "instructions": "Does this hit require escalation to a senior compliance officer or MLRO beyond standard resolution?",
    "criteria": {
      "true": "The hit involves a sanctions designation, terrorism financing, or proliferation — or the verdict is TRUE_MATCH regardless of category.",
      "false": "The hit can be resolved at analyst level: clear false positive, or a NEEDS_REVIEW that only requires additional data gathering."
    }
  },
  "additional_info_needed": {
    "type": "noul",
    "instructions": "Is additional identifying information required to reach a definitive verdict?",
    "criteria": {
      "true": "One or more critical biographic fields are unavailable and the verdict cannot be confirmed without them.",
      "false": "Available data is sufficient to reach a definitive TRUE_MATCH or FALSE_POSITIVE conclusion."
    }
  }
}
```

## Field Mapping from Source Data

| Template Field | JSON Source Path |
|---|---|
| `subject_name` | `data.summary.case_record.name` |
| `subject_dob` | `data.summary.case_record.date_of_birth` |
| `subject_nationality` | `data.summary.case_record.citizenship` |
| `subject_gender` | `data.summary.case_record.gender` |
| `subject_entity_type` | `data.summary.case_record.entity_type` |
| `hit_reference_id` | `data.summary.world_check[i].reference_id` |
| `hit_primary_name` | `data.summary.world_check[i].primary_name` |
| `hit_matched_term` | `data.summary.world_check[i].matched_term` |
| `hit_provider` | `data.summary.world_check[i].provider_type` |
| `hit_category` | `data.summary.world_check[i].key_data.category` |
| `hit_gender` | `data.summary.world_check[i].key_data.gender` |
| `hit_nationality` | `data.summary.world_check[i].key_data.location_details.nationality` |
| `hit_location` | `data.summary.world_check[i].key_data.addresses[*]` |
| `hit_aliases` | `data.summary.world_check[i].aliases[*]` |
| `comparison_fields` | `data.summary.world_check[i].comparison_data[*]` |
| `adverse_information` | `data.summary.world_check[i].further_information.details[*].text` |
| `associated_entities` | `data.summary.world_check[i].connections_and_relationships[*]` |
| `sources` | `data.summary.world_check[i].key_data.sources[*]` + `data.summary.world_check[i].sources[*]` |

## Original vs TypeSafe Mapping

| Original Prompt Field | TypeSafe Replacement |
|---|---|
| Free-text `name_analysis` | `name_match_strength` (score) + `name_distinguishing_power` (choice) |
| Free-text `biographic_comparison` | `dob_comparison` + `gender_comparison` + `nationality_comparison` + `location_comparison` (choice) + `biographic_conflict_count` (score) |
| Free-text `risk_assessment` | `risk_category_severity` (score) + `adverse_info_recency` (choice) + `associate_network_risk` (choice) |
| Mixed reasoning + verdict | Reasoning layer (10 typed questions) separated from Classification layer (4 typed questions) |
| `confidence: 0.0` float | `confidence_level` (score, 0-4 ordinal) |
| Free-text `recommended_action` | `escalation_required` (noul) + `additional_info_needed` (noul) |
| Free-text `key_factors` array | Eliminated — the reasoning questions themselves are the factors |
