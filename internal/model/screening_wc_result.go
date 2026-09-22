package model

import (
	"time"

	"github.com/go-viper/mapstructure/v2"
)

// ResolutionRiskEnum represents the risk level enum
type ResolutionRiskEnum string

const (
	ResolutionRiskEnumLow        ResolutionRiskEnum = "LOW"
	ResolutionRiskEnumMedium     ResolutionRiskEnum = "MEDIUM"
	ResolutionRiskEnumHigh       ResolutionRiskEnum = "HIGH"
	ResolutionRiskEnumProhibited ResolutionRiskEnum = "PROHIBITED"
)

// ResolutionStatusEnum represents the resolution status enum
type ResolutionStatusEnum string

const (
	ResolutionStatusEnumFalse ResolutionStatusEnum = "FALSE"
	ResolutionStatusEnumTrue  ResolutionStatusEnum = "TRUE"
)

// ScreeningResultCategory represents the screening result category
type ScreeningResultCategory string

// ScreeningResolveEnum represents the screening resolve status
type ScreeningResolveEnum string

// ScreeningWcResult represents a World-Check screening result
type ScreeningWcResult struct {
	*BaseEntity
	RiskLevel          ResolutionRiskEnum       `bson:"riskLevel" json:"riskLevel"`
	Category           ScreeningResultCategory  `bson:"category" json:"category"`
	AIRecommendation   map[string]interface{}   `bson:"aiRecommendation" json:"aiRecommendation"`
	Resolve            ScreeningResolveEnum     `bson:"resolve" json:"resolve"`
	ResolutionStatus   string                   `bson:"resolutionStatus" json:"resolutionStatus"`
	ResolutionReason   string                   `bson:"resolutionReason" json:"resolutionReason"`
	ResolutionRemark   string                   `bson:"resolutionRemark" json:"resolutionRemark"`
	ResultID           string                   `bson:"resultId" json:"resultId"`
	ReferenceID        string                   `bson:"referenceId" json:"referenceId"`
	MatchedTerm        string                   `bson:"matchedTerm" json:"matchedTerm"`
	PrimaryName        string                   `bson:"primaryName" json:"primaryName"`
	ReviewComment      string                   `bson:"reviewComment" json:"reviewComment"`
	Source             string                   `bson:"source" json:"source"`
	LastUpdated        *time.Time               `bson:"lastUpdated" json:"lastUpdated"`
	IsOgs              *bool                    `bson:"isOgs" json:"isOgs"`
	ComparisonData     []map[string]interface{} `bson:"comparisonData" json:"comparisonData"`
	Aliases            []map[string]interface{} `bson:"aliases" json:"aliases"`
	Sources            []map[string]interface{} `bson:"sources" json:"sources"`
	Keywords           []map[string]interface{} `bson:"keywords" json:"keywords"`
	ConnectionsAndRels []map[string]interface{} `bson:"connectionsAndRelationships" json:"connectionsAndRelationships"`
	RoleDetails        []map[string]interface{} `bson:"roleDetails" json:"roleDetails"`
	KeyData            map[string]interface{}   `bson:"keyData" json:"keyData"`
	FurtherInfo        map[string]interface{}   `bson:"furtherInformation" json:"furtherInformation"`
	Categories         []string                 `bson:"categories" json:"categories"`
	IsHighPriority     *bool                    `bson:"isHighPriority" json:"isHighPriority"`
	ReviewRequired     *bool                    `bson:"reviewRequired" json:"reviewRequired"`
	ReviewDate         *time.Time               `bson:"reviewDate" json:"reviewDate"`
	ReviewRequiredDate *time.Time               `bson:"reviewRequiredDate" json:"reviewRequiredDate"`
	ActionTypes        []string                 `bson:"actionTypes" json:"actionTypes"`
	ScreeningWc        *DBRef                   `bson:"screeningWc" json:"screeningWc"`
}

// GetRiskLevel returns the risk level with fallback logic
func (s *ScreeningWcResult) GetRiskLevel() ResolutionRiskEnum {
	if s.RiskLevel == "" && s.ResolutionStatus == string(ResolutionStatusEnumFalse) {
		return ResolutionRiskEnumLow
	}
	if s.ResolutionReason == string(ResolutionRiskEnumProhibited) {
		return ResolutionRiskEnumProhibited
	}
	return s.RiskLevel
}

// NewScreeningWcResultFromWorldCheckHits converts WorldCheckHits to ScreeningWcResult
func NewScreeningWcResultFromWorldCheckHits(w *WorldCheckHits) *ScreeningWcResult {
	result := &ScreeningWcResult{
		BaseEntity: &BaseEntity{
			PubID:     w.ReferenceID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ResultID:         w.ResultID,
		ReferenceID:      w.ReferenceID,
		MatchedTerm:      w.MatchedTerm,
		PrimaryName:      w.PrimaryName,
		ResolutionRemark: w.ResolutionRemark,
		LastUpdated:      convertToTimePtr(w.LastUpdated),
		Categories:       w.Categories,
		ReviewRequired:   &w.ReviewRequired,
		ActionTypes:      interfaceSliceToStringSlice(w.ActionTypes),
	}

	mapstructure.Decode(w, result)

	return result
}

func convertToTimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func interfaceSliceToStringSlice(items []interface{}) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
