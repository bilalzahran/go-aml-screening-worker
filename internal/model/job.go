package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobStatus string

const (
	Pending   JobStatus = "pending"
	Running   JobStatus = "running"
	Completed JobStatus = "completed"
	Failed    JobStatus = "failed"
)

type Job struct {
	ID         bson.ObjectID `json:"id" bson:"_id"`
	EventID    string        `json:"event_id" bson:"event_id"`
	Status     string        `json:"status" bson:"status"`
	ChildCount int           `json:"child_count" bson:"child_count"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}

type ChildJob struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	ParentID  bson.ObjectID `json:"parent_id" bson:"parent_id"`
	Status    string        `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
