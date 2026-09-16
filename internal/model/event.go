package model

import (
	"time"
)

type Event struct {
	Type             string           `json:"event_type" bson:"event_type"`
	ID               string           `json:"id" bson:"id"`
	TenantID         *string          `json:"tenant_id" bson:"tenant_id"`
	CaseID           string           `json:"case_id" bson:"case_id"`
	WorldCheck       []map[string]any `json:"world_check" bson:"world_check"`
	AIRecommendation []map[string]any `json:"ai_recommendation" bson:"ai_recommendation"`
	CreatedAt        time.Time
}
