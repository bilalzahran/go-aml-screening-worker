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

var _ repository.JobRepository = (*JobRepository)(nil)

type JobRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewJobRepository(database *mongo.Database, logger *slog.Logger) *JobRepository {
	return &JobRepository{
		collection: database.Collection("jobs"),
		logger:     logger.With("component", "mongo_job_repository"),
	}
}

func (r *JobRepository) Save(ctx context.Context, job *model.Job) error {
	filter := bson.D{{Key: "_id", Value: job.ID}}
	update := bson.D{{Key: "$set", Value: job}}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("upserting job %s: %w", job.ID, err)
	}

	r.logger.Debug("job saved", "job_id", job.ID)
	return nil
}
