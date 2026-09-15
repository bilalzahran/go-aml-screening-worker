package mongo

import (
	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ repository.ChildJobRepository = (*ChildJobRepository)(nil)

type ChildJobRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewChildJobRepository(collection *mongo.Collection, logger *slog.Logger) *ChildJobRepository {
	return &ChildJobRepository{
		collection: collection,
		logger:     logger,
	}
}

func (r *ChildJobRepository) Save(ctx context.Context, childJob *model.ChildJob) error {
	filter := bson.D{{Key: "_id", Value: childJob.ID}}
	update := bson.D{{Key: "$set", Value: childJob}}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("upserting child_job %s: %w", childJob.ID, err)
	}

	r.logger.Debug("child job saved", "child_job_id", childJob.ID)
	return nil
}
