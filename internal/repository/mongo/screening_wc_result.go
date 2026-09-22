package mongo

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
)

var _ repository.ScreeningWcResultRepository = (*ScreeningWcResultRepository)(nil)

type ScreeningWcResultRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewScreeningWcResultRepository(client *mongo.Client, database, collection string, logger *slog.Logger) *ScreeningWcResultRepository {
	return &ScreeningWcResultRepository{
		collection: client.Database(database).Collection(collection),
		logger:     logger.With("component", "mongo_screening_wc_result_repository"),
	}
}

func (r *ScreeningWcResultRepository) Save(ctx context.Context, result *model.ScreeningWcResult) error {
	filter := bson.D{{Key: "_id", Value: result.BaseEntity.ID}}
	update := bson.D{{Key: "$set", Value: result}}
	opts := options.UpdateOne().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("upserting screening_wc_result %s: %w", result.BaseEntity.ID, err)
	}

	r.logger.Debug("screening_wc_result saved", "id", result.BaseEntity.ID)
	return nil
}
