package model

// AIAnswer represents a single answer from the TypeSafe API
type AIAnswer struct {
	Type          string             `bson:"type" json:"type"`
	Choice        string             `bson:"choice,omitempty" json:"choice,omitempty"`
	Score         *float64           `bson:"score,omitempty" json:"score,omitempty"`
	Noul          *float64           `bson:"noul,omitempty" json:"noul,omitempty"`
	Legend        map[string]string  `bson:"legend,omitempty" json:"legend,omitempty"`
	Confidence    *float64           `bson:"confidence,omitempty" json:"confidence,omitempty"`
	Probabilities map[string]float64 `bson:"probabilities,omitempty" json:"probabilities,omitempty"`
}

// AIRecommendation contains reasoning and classification results from TypeSafe API
type AIRecommendation struct {
	Reasoning      map[string]AIAnswer `bson:"reasoning" json:"reasoning"`
	Classification map[string]AIAnswer `bson:"classification" json:"classification"`
}
