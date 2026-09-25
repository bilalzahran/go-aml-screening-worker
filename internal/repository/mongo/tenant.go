package mongo

import (
	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ repository.TenantRepository = (*TenantRepository)(nil)

type TenantRepository struct {
	collection *mongo.Collection
	logger     slog.Logger
}

func NewTenantRepository(client *mongo.Client, database, collection string, logger *slog.Logger) *TenantRepository {
	return &TenantRepository{
		collection: client.Database(database).Collection(collection),
		logger:     *logger,
	}
}

func (t *TenantRepository) FindByPubId(ctx context.Context, pubId string) (*model.Tenant, error) {

	filter := bson.M{"pubId": pubId}

	var tenant *model.Tenant

	err := t.collection.FindOne(ctx, filter).Decode(&tenant)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTenantNotFound
		}

		return nil, err
	}

	return tenant, nil
}
