package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
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
	*BaseEntity        `bson:",inline"`
	RiskLevel          ResolutionRiskEnum       `bson:"riskLevel" json:"riskLevel"`
	Category           ScreeningResultCategory  `bson:"category" json:"category"`
	AIRecommendation   *AIRecommendation        `bson:"aiRecommendation" json:"aiRecommendation"`
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
			PubID:     uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ResultID:              w.ResultID,
		ReferenceID:           w.ReferenceID,
		MatchedTerm:           w.MatchedTerm,
		PrimaryName:           w.PrimaryName,
		ResolutionRemark:      w.ResolutionRemark,
		ReviewComment:         interfaceToString(w.ReviewComment),
		Source:                w.ProviderType,
		LastUpdated:           convertToTimePtr(w.LastUpdated),
		Categories:            w.Categories,
		ReviewRequired:        &w.ReviewRequired,
		ActionTypes:           interfaceSliceToStringSlice(w.ActionTypes),
		ResolutionStatus:      interfaceToString(w.ResolutionStatus),
		ResolutionReason:      interfaceToString(w.ResolutionReason),
		ComparisonData:        structSliceToMapSlice(w.ComparisonData),
		Aliases:               structSliceToMapSlice(w.Aliases),
		Keywords:              structSliceToMapSlice(w.Keywords),
		Sources:               structSliceToMapSlice(w.KeyData.Sources),
		RoleDetails:           interfaceSliceToMapSlice(w.RoleDetails),
		ConnectionsAndRels:    interfaceSliceToMapSlice(w.ConnectionsAndRelationships),
		FurtherInfo:           structToMap(w.FurtherInformation),
		KeyData:               structToMap(w.KeyData),
	}

	// Handle RiskLevel conversion from interface{} to ResolutionRiskEnum
	if riskLevel, ok := w.RiskLevel.(string); ok {
		result.RiskLevel = ResolutionRiskEnum(riskLevel)
	}

	// Handle ReviewDate conversion from interface{} to *time.Time
	if reviewDate, ok := w.ReviewDate.(time.Time); ok {
		result.ReviewDate = &reviewDate
	}

	// Handle ReviewRequiredDate conversion if available
	if reviewRequiredDate, ok := w.ReviewRequiredDate.(time.Time); ok {
		result.ReviewRequiredDate = &reviewRequiredDate
	}

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

func interfaceToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func interfaceSliceToMapSlice(items []interface{}) []map[string]interface{} {
	if len(items) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result
}

func structSliceToMapSlice(items interface{}) []map[string]interface{} {
	if items == nil {
		return nil
	}

	// Marshal to JSON to respect json tags (snake_case)
	jsonData, err := json.Marshal(items)
	if err != nil {
		return nil
	}

	var decoded []map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		return nil
	}

	// Convert snake_case keys to camelCase for BSON
	return snakeToCamelCaseMapSlice(decoded)
}

func structToMap(item interface{}) map[string]interface{} {
	if item == nil {
		return nil
	}

	// Marshal to JSON to respect json tags (snake_case)
	jsonData, err := json.Marshal(item)
	if err != nil {
		return nil
	}

	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		return nil
	}

	// Convert snake_case keys to camelCase for BSON
	return snakeToCamelCaseMap(decoded)
}

// snakeToCamelCase converts snake_case string to camelCase
func snakeToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	result := parts[0]
	for _, part := range parts[1:] {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return result
}

// snakeToCamelCaseMap recursively converts all snake_case keys to camelCase
func snakeToCamelCaseMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		camelKey := snakeToCamelCase(k)
		if nestedMap, ok := v.(map[string]interface{}); ok {
			result[camelKey] = snakeToCamelCaseMap(nestedMap)
		} else if arr, ok := v.([]interface{}); ok {
			result[camelKey] = snakeToCamelCaseMapSliceInterface(arr)
		} else {
			result[camelKey] = v
		}
	}
	return result
}

// snakeToCamelCaseMapSlice converts all maps in a slice
func snakeToCamelCaseMapSlice(arr []map[string]interface{}) []map[string]interface{} {
	if len(arr) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(arr))
	for i, m := range arr {
		result[i] = snakeToCamelCaseMap(m)
	}
	return result
}

// snakeToCamelCaseMapSliceInterface converts maps inside interface{} slices
func snakeToCamelCaseMapSliceInterface(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result[i] = snakeToCamelCaseMap(m)
		} else {
			result[i] = item
		}
	}
	return result
}
