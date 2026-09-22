package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Event struct {
	Type             string           `json:"event_type" bson:"event_type"`
	ID               bson.ObjectID    `json:"id" bson:"_id"`
	Subject          OgsEventSubject  `json:"subject" bson:"subject"`
	TenantID         *string          `json:"tenant_id" bson:"tenant_id"`
	CaseID           string           `json:"case_id" bson:"case_id"`
	WorldCheck       []map[string]any `json:"world_check" bson:"world_check"`
	AIRecommendation []map[string]any `json:"ai_recommendation" bson:"ai_recommendation"`
	CreatedAt        time.Time        `json:"created_at" bson:"created_at"`
}

type OgsEventSubject struct {
	Name        string `json:"name"`
	DOB         string `json:"dob"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
	EntityType  string `json:"entity_type"`
}
