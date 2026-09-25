package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// BaseEntity is the base for all domain entities
type BaseEntity struct {
	ID             bson.ObjectID          `bson:"_id,omitempty" json:"-"`
	PubID          string                 `bson:"pubId" json:"pubId"`
	CreatedAt      time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time             `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	Attributes     map[string]interface{} `bson:"attributes" json:"attributes"`
	Version        int64                  `bson:"version" json:"version,omitempty"`
	CreatedBy      string                 `bson:"createdBy" json:"createdBy"`
	LastModifiedBy string                 `bson:"lastModifiedBy" json:"lastModifiedBy"`
}

func createBaseEntity() *BaseEntity {
	return &BaseEntity{
		PubID:     uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// DBRef represents a MongoDB database reference
type DBRef struct {
	Ref string        `bson:"$ref"`
	ID  bson.ObjectID `bson:"$id"`
	DB  string        `bson:"$db,omitempty"`
}
