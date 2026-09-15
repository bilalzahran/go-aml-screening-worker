package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        string          `json:"id"         bson:"_id"`
	Type      string          `json:"type"       bson:"type"`
	Source    string          `json:"source"     bson:"source"`
	Timestamp time.Time       `json:"timestamp"  bson:"timestamp"`
	Payload   json.RawMessage `json:"payload"    bson:"payload"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
}
