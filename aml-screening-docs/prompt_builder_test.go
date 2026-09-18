package screening

import (
	"fmt"
	"testing"
)

func TestBuildUserPrompt(t *testing.T) {
	subject := Subject{
		Name:        "王俊明",
		DOB:         "1970-07-05",
		Nationality: "CHN",
		Gender:      "MALE",
		EntityType:  "INDIVIDUAL",
	}

	hit := Hit{
		ReferenceID: "e_tr_wci_4995308",
		PrimaryName: "Junming WANG",
		MatchedTerm: "王俊明",
		Provider:    "WATCHLIST",
		Category:    "POST CONVICTION",
		Gender:      "MALE",
		Nationality: "CHN",
		Location:    "Shanxi, China",
		Aliases: []Alias{
			{Name: "Junming WANG", Type: "PRIMARY"},
			{Name: "王俊明", Type: "NATIVE_AKA", Language: "Chinese"},
			{Name: "WANG, Jun Ming", Type: "AKA"},
			{Name: "3769 0193 2494", Type: "NATIVE_AKA", Language: "Chinese — telegraph code"},
		},
		Comparisons: []ComparisonField{
			{Field: "NATIVE_AKA", Submitted: "王俊明", Matched: "王俊明", Result: "EXACT"},
			{Field: "DATE_OF_BIRTH", Submitted: "1970-07-05", Matched: "", Result: "N/A"},
			{Field: "NATIONALITY", Submitted: "CHN", Matched: "CHN", Result: "EXACT"},
			{Field: "GENDER", Submitted: "MALE", Matched: "MALE", Result: "EXACT"},
		},
		AdverseInfo: AdverseInfo{
			Categories: []string{"Special Interest Categories", "Other Bodies"},
			Keywords:   []string{"Narcotics Trafficking (M:1SK)", "Illegal Possession or Sale (M:1TF)"},
			Reports: []string{
				"May 2019 — sentenced to unspecified term of imprisonment of between 3 years, 8 months and 16 years imprisonment by Jinzhong Municipal Intermediate People's Court for narcotics trafficking.",
			},
		},
		Associates: []Associate{
			{Name: "Haigang ZHANG", Category: "POST CONVICTION", EntityType: "INDIVIDUAL"},
			{Name: "Bianli HOU", Category: "POST CONVICTION", EntityType: "INDIVIDUAL"},
			{Name: "Haina YANG", Category: "POST CONVICTION", EntityType: "INDIVIDUAL"},
		},
		Sources: []string{
			"http://www.sohu.com/a/316416193_127366?sec=wd",
			"http://www.legaldaily.com.cn/index/content/2019-05/24/content_7887167.htm",
		},
	}

	prompt := BuildUserPrompt(subject, hit)

	// Basic sanity checks
	if prompt == "" {
		t.Fatal("prompt should not be empty")
	}

	checks := []string{
		"王俊明",
		"1970-07-05",
		"e_tr_wci_4995308",
		"Junming WANG",
		"NATIVE_AKA",
		"Narcotics Trafficking",
		"POST CONVICTION",
		"Respond ONLY with JSON",
	}
	for _, want := range checks {
		if !containsStr(prompt, want) {
			t.Errorf("prompt missing expected substring: %q", want)
		}
	}

	fmt.Println(prompt)
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && findStr(s, substr)
}

func findStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
