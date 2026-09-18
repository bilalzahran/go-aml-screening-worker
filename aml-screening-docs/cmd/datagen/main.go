package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// ---------- data pools ----------

type nameEntry struct {
	Chinese    string
	Romanized  string
	Telegraph  string
	Gender     string
}

var namePool = []nameEntry{
	{"王俊明", "Junming WANG", "3769 0193 2494", "MALE"},
	{"陈志强", "Zhiqiang CHEN", "7115 1807 1730", "MALE"},
	{"刘建国", "Jianguo LIU", "0491 1696 0948", "MALE"},
	{"黄美玲", "Meiling HUANG", "7806 5765 3781", "FEMALE"},
	{"林秀英", "Xiuying LIN", "2651 4423 5765", "FEMALE"},
	{"张伟", "Wei ZHANG", "1728 0251", "MALE"},
	{"李娜", "Na LI", "2621 1282", "FEMALE"},
	{"赵小龙", "Xiaolong ZHAO", "6392 1420 7893", "MALE"},
	{"杨丽华", "Lihua YANG", "2799 4409 5478", "FEMALE"},
	{"吴国栋", "Guodong WU", "0702 0948 2637", "MALE"},
	{"周文杰", "Wenjie ZHOU", "0719 2429 2638", "MALE"},
	{"孙雪梅", "Xuemei SUN", "1327 7185 2734", "FEMALE"},
	{"徐天宇", "Tianyu XU", "1776 1131 1342", "MALE"},
	{"马金凤", "Jinfeng MA", "7456 6855 7685", "FEMALE"},
	{"朱海波", "Haibo ZHU", "2612 3189 3134", "MALE"},
	{"何晓燕", "Xiaoyan HE", "0149 1420 3601", "FEMALE"},
	{"郭明亮", "Mingliang GUO", "6753 2494 0081", "MALE"},
	{"谢玉兰", "Yulan XIE", "6200 3768 5765", "FEMALE"},
	{"邓文龙", "Wenlong DENG", "6772 2429 7893", "MALE"},
	{"曹慧敏", "Huimin CAO", "2580 1979 2404", "FEMALE"},
	{"宋建华", "Jianhua SONG", "1345 1696 5478", "MALE"},
	{"韩志勇", "Zhiyong HAN", "7281 1807 0516", "MALE"},
	{"冯丽芳", "Lifang FENG", "7458 4409 5364", "FEMALE"},
	{"彭国平", "Guoping PENG", "1756 0948 1627", "MALE"},
	{"蒋玉梅", "Yumei JIANG", "5592 3768 2734", "FEMALE"},
	{"沈伟民", "Weimin SHEN", "3088 0251 3046", "MALE"},
	{"唐晓东", "Xiaodong TANG", "0781 1420 2639", "MALE"},
	{"许文静", "Wenjing XU", "6079 2429 7234", "FEMALE"},
	{"叶志明", "Zhiming YE", "0673 1807 2494", "MALE"},
	{"余建林", "Jianlin YU", "0151 1696 2651", "MALE"},
}

type sourceEntry struct {
	Abbreviation string
	Name         string
	Description  string
}

var sanctionSources = []sourceEntry{
	{"OFAC", "OFAC - SDN List", "US Treasury OFAC Specially Designated Nationals list."},
	{"EU-SL", "EU Consolidated Sanctions List", "European Union consolidated list of designated persons and entities."},
	{"UN-SC", "UN Security Council Consolidated List", "United Nations Security Council consolidated list."},
	{"MAS-TF", "MAS Terrorism Financing List", "Monetary Authority of Singapore designated persons."},
	{"HKMA-SF", "HKMA Sanctions & Freezing", "Hong Kong Monetary Authority sanctions orders."},
}

var pepSources = []sourceEntry{
	{"PEP-CN", "China PEP Database", "Politically exposed persons in China."},
	{"PEP-ASEAN", "ASEAN PEP Registry", "Politically exposed persons in ASEAN member states."},
}

var adverseMediaSources = []sourceEntry{
	{"AM-GLOBAL", "Global Adverse Media", "Negative news and adverse media screening."},
	{"AM-APAC", "APAC Adverse Media", "Asia-Pacific focused adverse media monitoring."},
}

var categories = []string{"SANCTIONS", "PEP", "ADVERSE MEDIA", "POST CONVICTION", "SPECIAL INTEREST"}

type countryEntry struct {
	Code string
	Name string
}

var countries = []countryEntry{
	{"CHN", "China"},
	{"MYS", "Malaysia"},
	{"SGP", "Singapore"},
	{"IDN", "Indonesia"},
	{"THA", "Thailand"},
	{"HKG", "Hong Kong"},
	{"TWN", "Taiwan"},
	{"VNM", "Vietnam"},
	{"PHL", "Philippines"},
	{"KHM", "Cambodia"},
}

var cities = map[string][]string{
	"CHN": {"Beijing", "Shanghai", "Guangzhou", "Shenzhen", "Chengdu", "Hangzhou", "Wuhan"},
	"MYS": {"Kuala Lumpur", "Penang", "Johor Bahru", "Kota Kinabalu"},
	"SGP": {"Singapore"},
	"IDN": {"Jakarta", "Surabaya", "Bandung", "Medan"},
	"THA": {"Bangkok", "Chiang Mai", "Phuket"},
	"HKG": {"Hong Kong"},
	"TWN": {"Taipei", "Kaohsiung", "Taichung"},
	"VNM": {"Hanoi", "Ho Chi Minh City", "Da Nang"},
	"PHL": {"Manila", "Cebu", "Davao"},
	"KHM": {"Phnom Penh", "Siem Reap"},
}

var adverseDetails = []string{
	"Under investigation by local financial authority for suspected money laundering activities.",
	"Named in media reports related to corruption investigation.",
	"Subject of regulatory inquiry regarding suspicious transaction patterns.",
	"Mentioned in leaked financial documents related to offshore accounts.",
	"Referenced in court filings related to bribery allegations.",
	"Listed in investigative journalism report on illicit financial networks.",
	"Named as associate of designated individual in enforcement action.",
	"Subject of parliamentary inquiry regarding misuse of public funds.",
}

var sanctionDetails = []string{
	"Designated under counter-terrorism financing regulations.",
	"Added to sanctions list for involvement in proliferation financing.",
	"Designated for facilitating illicit financial flows.",
	"Listed for connection to weapons trafficking network.",
	"Sanctioned for involvement in drug trafficking operations.",
}

var pepRoles = []string{
	"Member of Provincial People's Congress",
	"Deputy Director, Municipal Planning Bureau",
	"Vice Chairman, State-Owned Enterprise",
	"Director General, Trade Regulatory Authority",
	"Senior Adviser, Ministry of Finance",
}

// ---------- types ----------

type TestCase struct {
	Comment        string `json:"_comment"`
	ExpectedVerdict string `json:"expected_verdict"`
	TestID         string `json:"test_id"`
	Data           any    `json:"data"`
}

type Difficulty string

const (
	ObviousFalsePositive Difficulty = "OBVIOUS_FALSE_POSITIVE"
	ObviousTrueMatch     Difficulty = "OBVIOUS_TRUE_MATCH"
	Ambiguous            Difficulty = "AMBIGUOUS"
)

// ---------- generator ----------

type generator struct {
	rng *rand.Rand
}

func newGenerator(seed int64) *generator {
	return &generator{rng: rand.New(rand.NewSource(seed))}
}

func (g *generator) pick(n int) int {
	return g.rng.Intn(n)
}

func (g *generator) pickName() nameEntry {
	return namePool[g.pick(len(namePool))]
}

func (g *generator) pickCountry() countryEntry {
	return countries[g.pick(len(countries))]
}

func (g *generator) pickCity(code string) string {
	c := cities[code]
	if len(c) == 0 {
		return "Unknown"
	}
	return c[g.pick(len(c))]
}

func (g *generator) pickCategory() string {
	return categories[g.pick(len(categories))]
}

func (g *generator) pickSource(category string) sourceEntry {
	switch category {
	case "SANCTIONS":
		return sanctionSources[g.pick(len(sanctionSources))]
	case "PEP":
		return pepSources[g.pick(len(pepSources))]
	default:
		return adverseMediaSources[g.pick(len(adverseMediaSources))]
	}
}

func (g *generator) randomDOB(baseYear int, spreadYears int) string {
	year := baseYear - spreadYears + g.pick(2*spreadYears+1)
	month := 1 + g.pick(12)
	day := 1 + g.pick(28)
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func (g *generator) generateHit(
	subjectName nameEntry,
	subjectDOB string,
	subjectCountry countryEntry,
	hitIndex int,
	casePrefix string,
	difficulty Difficulty,
) map[string]any {
	refID := fmt.Sprintf("e_tr_wci_%s_%04d", casePrefix, hitIndex)
	resultID := fmt.Sprintf("%s-result-%04d", casePrefix, hitIndex)

	var hitName nameEntry
	var hitDOB string
	var hitCountry countryEntry
	var hitGender string
	var category string
	var matchedTerm string

	switch difficulty {
	case ObviousFalsePositive:
		// Different person: different gender, large DOB gap, different country
		for {
			hitName = g.pickName()
			if hitName.Gender != subjectName.Gender {
				break
			}
		}
		// DOB 15-30 years off
		hitDOB = g.randomDOB(1960, 15)
		for {
			hitCountry = g.pickCountry()
			if hitCountry.Code != subjectCountry.Code {
				break
			}
		}
		hitGender = hitName.Gender
		category = categories[g.pick(len(categories))]
		matchedTerm = subjectName.Chinese // name matched but everything else conflicts

	case ObviousTrueMatch:
		// Exact match: same name, same DOB, same gender, same country, sanctions
		hitName = subjectName
		hitDOB = subjectDOB
		hitCountry = subjectCountry
		hitGender = subjectName.Gender
		category = "SANCTIONS"
		matchedTerm = subjectName.Chinese

	case Ambiguous:
		// Partial match: similar name, some fields match, some missing
		ambiguityType := g.pick(5)
		switch ambiguityType {
		case 0: // Same name, DOB off by 1-2 years, same country
			hitName = subjectName
			hitDOB = g.randomDOB(1985, 2)
			hitCountry = subjectCountry
			hitGender = subjectName.Gender
			category = g.pickCategory()
			matchedTerm = subjectName.Chinese
		case 1: // Same name, all biographic fields null/missing
			hitName = subjectName
			hitDOB = ""
			hitCountry = subjectCountry
			hitGender = ""
			category = "PEP"
			matchedTerm = subjectName.Chinese
		case 2: // Different romanization, same Chinese characters
			hitName = subjectName
			hitDOB = subjectDOB
			hitCountry = g.pickCountry()
			hitGender = subjectName.Gender
			category = g.pickCategory()
			matchedTerm = subjectName.Romanized
		case 3: // Telegraph code match, biographic partial
			hitName = subjectName
			hitDOB = ""
			hitCountry = countryEntry{Code: "", Name: ""}
			hitGender = subjectName.Gender
			category = g.pickCategory()
			matchedTerm = subjectName.Telegraph
		default: // Name is close but not exact (different person from pool with same gender)
			for {
				hitName = g.pickName()
				if hitName.Chinese != subjectName.Chinese && hitName.Gender == subjectName.Gender {
					break
				}
			}
			hitDOB = subjectDOB
			hitCountry = subjectCountry
			hitGender = subjectName.Gender
			category = "ADVERSE MEDIA"
			matchedTerm = subjectName.Chinese
		}
	}

	src := g.pickSource(category)
	hitCity := g.pickCity(hitCountry.Code)

	// Build comparison_data
	comparisonData := []map[string]any{
		{"typeId": "", "name": "NATIVE_AKA", "submittedValue": subjectName.Chinese, "matchedValue": hitName.Chinese},
	}
	if hitDOB != "" {
		comparisonData = append(comparisonData, map[string]any{
			"typeId": "SFCT_2", "name": "DATE_OF_BIRTH",
			"submittedValue": subjectDOB, "matchedValue": hitDOB,
		})
	}
	if hitCountry.Code != "" {
		comparisonData = append(comparisonData, map[string]any{
			"typeId": "SFCT_5", "name": "NATIONALITY",
			"submittedValue": subjectCountry.Code, "matchedValue": hitCountry.Code,
		})
	}
	if hitGender != "" {
		comparisonData = append(comparisonData, map[string]any{
			"typeId": "SFCT_1", "name": "GENDER",
			"submittedValue": subjectName.Gender, "matchedValue": hitGender,
		})
	}

	// Build aliases
	aliases := []map[string]any{
		{"type": "PRIMARY", "orginal_script": hitName.Romanized, "submitted_data": subjectName.Chinese, "matched_data": hitName.Romanized, "language": ""},
		{"type": "NATIVE_AKA", "orginal_script": hitName.Chinese, "submitted_data": subjectName.Chinese, "matched_data": hitName.Chinese, "language": "Chinese"},
	}
	if hitName.Telegraph != "" && g.pick(2) == 0 {
		aliases = append(aliases, map[string]any{
			"type": "NATIVE_AKA", "orginal_script": hitName.Telegraph,
			"submitted_data": subjectName.Chinese, "matched_data": hitName.Telegraph, "language": "Chinese",
		})
	}

	// Build address
	var addressCountryObj map[string]any
	if hitCountry.Code != "" {
		addressCountryObj = map[string]any{"code": hitCountry.Code, "name": hitCountry.Name}
	} else {
		addressCountryObj = map[string]any{"code": nil, "name": nil}
	}

	// Build further_information
	var detailText string
	switch {
	case category == "SANCTIONS":
		detailText = sanctionDetails[g.pick(len(sanctionDetails))]
	case category == "PEP":
		detailText = fmt.Sprintf("Currently serving as %s.", pepRoles[g.pick(len(pepRoles))])
	default:
		detailText = adverseDetails[g.pick(len(adverseDetails))]
	}

	// Build location details
	locationDetails := map[string]any{"location": hitCountry.Name}
	if hitCountry.Code == "" {
		locationDetails["location"] = nil
	}

	// Determine gender value for key_data (may be null for ambiguous)
	var genderVal any = hitGender
	if hitGender == "" {
		genderVal = nil
	}

	// Determine DOB-related fields
	var dobVal any = hitDOB
	if hitDOB == "" {
		dobVal = nil
	}
	_ = dobVal // DOB is in comparison_data, not top-level

	// PEP status
	pepStatus := ""
	if category == "PEP" {
		pepStatus = "CURRENT"
	}

	// Categories mapping
	var catList []string
	switch category {
	case "SANCTIONS":
		catList = []string{"Sanctions"}
	case "PEP":
		catList = []string{"PEP"}
	case "ADVERSE MEDIA":
		catList = []string{"Adverse Media"}
	case "POST CONVICTION":
		catList = []string{"Special Interest Categories"}
	default:
		catList = []string{"Special Interest Categories"}
	}

	lastUpdated := time.Date(2026, 9, 15, 10, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
	lastUpdated = "2026-09-15T10:00:00.000Z"

	hit := map[string]any{
		"reference_id":      refID,
		"result_id":         resultID,
		"matched_term":      matchedTerm,
		"primary_name":      hitName.Romanized,
		"provider_type":     "WATCHLIST",
		"resolution_status": nil,
		"resolution_reason": nil,
		"risk_level":        nil,
		"resolution_remark": "",
		"review_comment":    nil,
		"review_required":   false,
		"review_required_date": nil,
		"review_date":       nil,
		"last_updated":      lastUpdated,
		"pep_status":        pepStatus,
		"comparison_data":   comparisonData,
		"aliases":           aliases,
		"key_data": map[string]any{
			"gender":   genderVal,
			"category": category,
			"sources":  []map[string]any{{"abbreviation": src.Abbreviation, "name": src.Name, "description": src.Description}},
			"previous_roles": []any{},
			"addresses": []map[string]any{{
				"city":     hitCity,
				"country":  addressCountryObj,
				"postCode": nil,
				"region":   nil,
				"street":   nil,
			}},
			"identity_documents": []any{},
			"location_details":   locationDetails,
		},
		"further_information": map[string]any{
			"details": []map[string]any{{
				"detailType": "REPORTS",
				"text":       detailText,
				"title":      "REPORTS",
			}},
		},
		"keywords": []map[string]any{{
			"abbreviation": src.Abbreviation,
			"name":         src.Name,
			"description":  src.Description,
		}},
		"sources": []map[string]any{{
			"url":     fmt.Sprintf("https://example.com/%s-%s-%04d", casePrefix, category, hitIndex),
			"caption": nil,
		}},
		"connections_and_relationships": []any{},
		"role_details":                  []any{},
		"update_category":              "UNKNOWN",
		"categories":                   catList,
		"action_types":                 []any{},
		"_difficulty":                  string(difficulty),
	}

	return hit
}

func (g *generator) generateTestCase(
	testID string,
	comment string,
	expectedVerdict string,
	casePrefix string,
	totalHits int,
	difficultyDist map[Difficulty]float64,
) TestCase {
	subject := g.pickName()
	subjectDOB := g.randomDOB(1980, 10)
	subjectCountry := g.pickCountry()

	// Build distribution list
	type hitSpec struct {
		diff Difficulty
	}
	var specs []hitSpec

	ambiguousCount := int(float64(totalHits) * difficultyDist[Ambiguous])
	trueMatchCount := int(float64(totalHits) * difficultyDist[ObviousTrueMatch])
	falsePositiveCount := totalHits - ambiguousCount - trueMatchCount

	for i := 0; i < ambiguousCount; i++ {
		specs = append(specs, hitSpec{Ambiguous})
	}
	for i := 0; i < trueMatchCount; i++ {
		specs = append(specs, hitSpec{ObviousTrueMatch})
	}
	for i := 0; i < falsePositiveCount; i++ {
		specs = append(specs, hitSpec{ObviousFalsePositive})
	}

	// Shuffle
	g.rng.Shuffle(len(specs), func(i, j int) { specs[i], specs[j] = specs[j], specs[i] })

	hits := make([]map[string]any, len(specs))
	for i, spec := range specs {
		hits[i] = g.generateHit(subject, subjectDOB, subjectCountry, i+1, casePrefix, spec.diff)
	}

	data := map[string]any{
		"summary": map[string]any{
			"case_record": map[string]any{
				"name":                  subject.Chinese,
				"case_id":              fmt.Sprintf("%s-case-id", casePrefix),
				"case_system_id":       fmt.Sprintf("%s-case-system-id", casePrefix),
				"group_name":           "",
				"group_id":            fmt.Sprintf("%s-group-id", casePrefix),
				"ongoing_screening":    "ONGOING",
				"entity_type":         "INDIVIDUAL",
				"status":              "UNARCHIVED",
				"name_transposition":   false,
				"outstanding_actions":  true,
				"gender":              subject.Gender,
				"date_of_birth":       subjectDOB,
				"citizenship":         subjectCountry.Code,
				"country_location":    subjectCountry.Name,
				"place_of_birth":      g.pickCity(subjectCountry.Code),
				"registration_country": "",
				"identification_numbers": "",
			},
			"key_findings": map[string]any{
				"world_check": totalHits,
				"watchlist":   0,
			},
			"world_check": hits,
		},
	}

	return TestCase{
		Comment:         comment,
		ExpectedVerdict: expectedVerdict,
		TestID:          testID,
		Data:            data,
	}
}

func main() {
	outPath := flag.String("out", "test_screening_data.json", "path to test data JSON file")
	seed := flag.Int64("seed", 42, "random seed for reproducible generation")
	flag.Parse()

	// Read existing test cases
	raw, err := os.ReadFile(*outPath)
	if err != nil {
		log.Fatalf("failed to read existing test data: %v", err)
	}

	var existing []json.RawMessage
	if err := json.Unmarshal(raw, &existing); err != nil {
		log.Fatalf("failed to parse existing test data: %v", err)
	}

	// Filter out any previously generated TC-016/017/018
	var filtered []json.RawMessage
	for _, entry := range existing {
		var peek struct {
			TestID string `json:"test_id"`
		}
		_ = json.Unmarshal(entry, &peek)
		switch peek.TestID {
		case "TC-016-BULK-ALL-AMBIGUOUS", "TC-017-BULK-HALF-AMBIGUOUS", "TC-018-BULK-MOSTLY-OBVIOUS":
			log.Printf("removing existing %s (will regenerate)", peek.TestID)
		default:
			filtered = append(filtered, entry)
		}
	}

	log.Printf("kept %d existing test cases", len(filtered))

	gen := newGenerator(*seed)

	// TC-016: 50 hits, 100% ambiguous
	tc016 := gen.generateTestCase(
		"TC-016-BULK-ALL-AMBIGUOUS",
		"TEST CASE 16: BULK — 50 hits, ALL ambiguous (100% need LLM). Stress-tests fan-out with no pre-filter shortcut.",
		"MIXED",
		"tc016",
		50,
		map[Difficulty]float64{Ambiguous: 1.0, ObviousTrueMatch: 0, ObviousFalsePositive: 0},
	)

	// TC-017: 100 hits, ~50% ambiguous
	tc017 := gen.generateTestCase(
		"TC-017-BULK-HALF-AMBIGUOUS",
		"TEST CASE 17: BULK — 100 hits, ~50% ambiguous (half need LLM, half obvious). Tests mixed pre-filter + fan-out.",
		"MIXED",
		"tc017",
		100,
		map[Difficulty]float64{Ambiguous: 0.5, ObviousTrueMatch: 0.15, ObviousFalsePositive: 0.35},
	)

	// TC-018: 200 hits, ~10% ambiguous
	tc018 := gen.generateTestCase(
		"TC-018-BULK-MOSTLY-OBVIOUS",
		"TEST CASE 18: BULK — 200 hits, ~10% ambiguous (mostly obvious, tests pre-filter effectiveness). Only ~20 hits need LLM.",
		"MIXED",
		"tc018",
		200,
		map[Difficulty]float64{Ambiguous: 0.1, ObviousTrueMatch: 0.05, ObviousFalsePositive: 0.85},
	)

	// Marshal new test cases
	newCases := []TestCase{tc016, tc017, tc018}
	for _, tc := range newCases {
		b, err := json.Marshal(tc)
		if err != nil {
			log.Fatalf("failed to marshal %s: %v", tc.TestID, err)
		}
		filtered = append(filtered, json.RawMessage(b))
	}

	// Write combined output
	out, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal final output: %v", err)
	}

	if err := os.WriteFile(*outPath, append(out, '\n'), 0644); err != nil {
		log.Fatalf("failed to write output: %v", err)
	}

	// Summary
	for _, tc := range newCases {
		data := tc.Data.(map[string]any)
		summary := data["summary"].(map[string]any)
		hits := summary["world_check"].([]map[string]any)

		counts := map[Difficulty]int{}
		for _, h := range hits {
			d := Difficulty(h["_difficulty"].(string))
			counts[d]++
		}

		log.Printf("generated %s: %d hits (ambiguous=%d, obvious_true=%d, obvious_false=%d)",
			tc.TestID, len(hits),
			counts[Ambiguous], counts[ObviousTrueMatch], counts[ObviousFalsePositive])
	}

	log.Printf("wrote %d total test cases to %s", len(filtered), *outPath)
}
