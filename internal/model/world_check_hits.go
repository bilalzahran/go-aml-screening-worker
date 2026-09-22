package model

import (
	"time"

	"github.com/go-viper/mapstructure/v2"
)

// MapToWorldCheckHits converts a map[string]any to WorldCheckHits
func MapToWorldCheckHits(data map[string]any) (*WorldCheckHits, error) {
	var result WorldCheckHits

	err := mapstructure.Decode(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type WorldCheckHits struct {
	RiskLevel                   interface{}        `json:"risk_level"`
	ActionTypes                 []interface{}      `json:"action_types"`
	ComparisonData              []ComparisonData   `json:"comparison_data"`
	KeyData                     KeyData            `json:"key_data"`
	Keywords                    []Keyword          `json:"keywords"`
	LastUpdated                 time.Time          `json:"last_updated"`
	PepStatus                   string             `json:"pep_status"`
	PrimaryName                 string             `json:"primary_name"`
	ResolutionStatus            interface{}        `json:"resolution_status"`
	ConnectionsAndRelationships []interface{}      `json:"connections_and_relationships"`
	ResolutionReason            interface{}        `json:"resolution_reason"`
	ReviewComment               interface{}        `json:"review_comment"`
	ReviewRequiredDate          interface{}        `json:"review_required_date"`
	UpdateCategory              string             `json:"update_category"`
	Aliases                     []Alias            `json:"aliases"`
	ReferenceID                 string             `json:"reference_id"`
	ReviewRequired              bool               `json:"review_required"`
	RoleDetails                 []interface{}      `json:"role_details"`
	Sources                     []interface{}      `json:"sources"`
	Categories                  []string           `json:"categories"`
	FurtherInformation          FurtherInformation `json:"further_information"`
	MatchedTerm                 string             `json:"matched_term"`
	ProviderType                string             `json:"provider_type"`
	ResolutionRemark            string             `json:"resolution_remark"`
	ResultID                    string             `json:"result_id"`
	ReviewDate                  interface{}        `json:"review_date"`
}


type ComparisonData struct {
	MatchedValue   string `json:"matchedValue"`
	Name           string `json:"name"`
	SubmittedValue string `json:"submittedValue"`
	TypeID         string `json:"typeId"`
}

type KeyData struct {
	Addresses         []interface{}  `json:"addresses"`
	Category          string         `json:"category"`
	Gender            string         `json:"gender"`
	IdentityDocuments []interface{}  `json:"identity_documents"`
	LocationDetails   LocationDetail `json:"location_details"`
	PreviousRoles     []interface{}  `json:"previous_roles"`
	Sources           []Source       `json:"sources"`
}

type LocationDetail struct {
	Nationality string `json:"nationality"`
	Location    string `json:"location"`
}

type Source struct {
	Abbreviation string `json:"abbreviation"`
	Description  string `json:"description"`
	Name         string `json:"name"`
}

type Keyword struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Description  string `json:"description"`
}

type Alias struct {
	SubmittedData string `json:"submitted_data"`
	Type          string `json:"type"`
	Language      string `json:"language"`
	MatchedData   string `json:"matched_data"`
	OrginalScript string `json:"orginal_script"`
}

type FurtherInformation struct {
	Details []FurtherInformationDetail `json:"details"`
}

type FurtherInformationDetail struct {
	Title      string `json:"title"`
	DetailType string `json:"detailType"`
	Text       string `json:"text"`
}
