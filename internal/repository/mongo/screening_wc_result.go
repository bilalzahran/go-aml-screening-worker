package mongo

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
)

var _ repository.ScreeningWcResultRepository = (*ScreeningWcResultRepository)(nil)

type ScreeningWcResultRepository struct {
	tenantResolver repository.TenantDbResolver
	logger         *slog.Logger
}

func NewScreeningWcResultRepository(tenantResolver repository.TenantDbResolver, logger *slog.Logger) *ScreeningWcResultRepository {
	return &ScreeningWcResultRepository{
		tenantResolver: tenantResolver,
		logger:         logger.With("component", "mongo_screening_wc_result_repository"),
	}
}

func (r *ScreeningWcResultRepository) Save(ctx context.Context, tenantId string, result *model.ScreeningWcResult) error {
	filter := bson.D{{Key: "_id", Value: result.BaseEntity.ID}}
	update := bson.D{{Key: "$set", Value: result}}
	opts := options.UpdateOne().SetUpsert(true)

	database, err := r.tenantResolver.Database(ctx, tenantId)
	if err != nil {
		return fmt.Errorf("error when resolving db: %w", err)
	}
	collection := database.Collection("screening_wc_results")

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("upserting screening_wc_result %s: %w", result.BaseEntity.ID, err)
	}

	r.logger.Debug("screening_wc_result saved", "id", result.BaseEntity.ID)
	return nil
}
