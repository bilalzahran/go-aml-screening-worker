package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TenantDbResolver interface {
	Database(ctx context.Context, tenantPubId string) (*mongo.Database, error)
}
