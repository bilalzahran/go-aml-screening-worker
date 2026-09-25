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

var _ repository.EventRepository = (*EventRepository)(nil)

type EventRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewEventRepository(database *mongo.Database, logger *slog.Logger) *EventRepository {
	return &EventRepository{
		collection: database.Collection("events"),
		logger:     logger.With("component", "mongo_repository"),
	}
}

func (r *EventRepository) Save(ctx context.Context, event *model.Event) error {
	filter := bson.D{{Key: "_id", Value: event.ID}}
	update := bson.D{{Key: "$set", Value: event}}
	opts := options.UpdateOne().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("upserting event %s: %w", event.ID, err)
	}

	r.logger.Debug("event saved", "event_id", event.ID, "event_type", event.Type)
	return nil
}
