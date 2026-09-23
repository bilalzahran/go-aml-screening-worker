package model

import (
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

// MapToWorldCheckHits converts a map[string]any to WorldCheckHits
func MapToWorldCheckHits(data map[string]any) (*WorldCheckHits, error) {
	var result WorldCheckHits

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &result,
		TagName:  "json", // Use json tags for field mapping
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			stringToTimeHook,
		),
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// stringToTimeHook converts string timestamps to time.Time
func stringToTimeHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	switch v := data.(type) {
	case string:
		if v == "" {
			return time.Time{}, nil
		}
		return time.Parse(time.RFC3339Nano, v)
	case time.Time:
		return v, nil
	}

	return data, nil
}

type WorldCheckHits struct {
	RiskLevel                   any                `json:"risk_level" bson:"riskLevel"`
	ActionTypes                 []any              `json:"action_types" bson:"actionTypes"`
	ComparisonData              []ComparisonData   `json:"comparison_data" bson:"comparisonData"`
	KeyData                     KeyData            `json:"key_data" bson:"keyData"`
	Keywords                    []Keyword          `json:"keywords" bson:"keywords"`
	LastUpdated                 time.Time          `json:"last_updated" bson:"lastUpdated"`
	PepStatus                   string             `json:"pep_status" bson:"pepStatus"`
	PrimaryName                 string             `json:"primary_name" bson:"primaryName"`
	ResolutionStatus            any                `json:"resolution_status" bson:"resolutionStatus"`
	ConnectionsAndRelationships []any              `json:"connections_and_relationships" bson:"connectionsAndRelationships"`
	ResolutionReason            any                `json:"resolution_reason" bson:"resolutionReason"`
	ReviewComment               any                `json:"review_comment" bson:"reviewComment"`
	ReviewRequiredDate          any                `json:"review_required_date" bson:"reviewRequiredDate"`
	UpdateCategory              string             `json:"update_category" bson:"updateCategory"`
	Aliases                     []Alias            `json:"aliases" bson:"aliases"`
	ReferenceID                 string             `json:"reference_id" bson:"referenceId"`
	ReviewRequired              bool               `json:"review_required" bson:"reviewRequired"`
	RoleDetails                 []any              `json:"role_details" bson:"roleDetails"`
	Sources                     []any              `json:"sources" bson:"sources"`
	Categories                  []string           `json:"categories" bson:"categories"`
	FurtherInformation          FurtherInformation `json:"further_information" bson:"furtherInformation"`
	MatchedTerm                 string             `json:"matched_term" bson:"matchedTerm"`
	ProviderType                string             `json:"provider_type" bson:"providerType"`
	ResolutionRemark            string             `json:"resolution_remark" bson:"resolutionRemark"`
	ResultID                    string             `json:"result_id" bson:"resultId"`
	ReviewDate                  any                `json:"review_date" bson:"reviewDate"`
}


type ComparisonData struct {
	MatchedValue   string `json:"matchedValue" bson:"matchedValue"`
	Name           string `json:"name" bson:"name"`
	SubmittedValue string `json:"submittedValue" bson:"submittedValue"`
	TypeID         string `json:"typeId" bson:"typeId"`
}

type KeyData struct {
	Addresses         []any          `json:"addresses" bson:"addresses"`
	Category          string         `json:"category" bson:"category"`
	Gender            string         `json:"gender" bson:"gender"`
	IdentityDocuments []any          `json:"identity_documents" bson:"identityDocuments"`
	LocationDetails   LocationDetail `json:"location_details" bson:"locationDetails"`
	PreviousRoles     []any          `json:"previous_roles" bson:"previousRoles"`
	Sources           []Source       `json:"sources" bson:"sources"`
}

type LocationDetail struct {
	Nationality string `json:"nationality" bson:"nationality"`
	Location    string `json:"location" bson:"location"`
}

type Source struct {
	Abbreviation string `json:"abbreviation" bson:"abbreviation"`
	Description  string `json:"description" bson:"description"`
	Name         string `json:"name" bson:"name"`
}

type Keyword struct {
	Name         string `json:"name" bson:"name"`
	Abbreviation string `json:"abbreviation" bson:"abbreviation"`
	Description  string `json:"description" bson:"description"`
}

type Alias struct {
	SubmittedData string `json:"submitted_data" bson:"submittedData"`
	Type          string `json:"type" bson:"type"`
	Language      string `json:"language" bson:"language"`
	MatchedData   string `json:"matched_data" bson:"matchedData"`
	OrginalScript string `json:"orginal_script" bson:"orginalScript"`
}

type FurtherInformation struct {
	Details []FurtherInformationDetail `json:"details" bson:"details"`
}

type FurtherInformationDetail struct {
	Title      string `json:"title" bson:"title"`
	DetailType string `json:"detailType" bson:"detailType"`
	Text       string `json:"text" bson:"text"`
}
