package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
)

type EventService struct {
	repo   repository.EventRepository
	logger *slog.Logger
}

func NewEventService(repo repository.EventRepository, logger *slog.Logger) *EventService {
	return &EventService{
		repo:   repo,
		logger: logger.With("component", "event_service"),
	}
}

func (s *EventService) Process(ctx context.Context, event *model.Event) error {
	if event.ID.IsZero() {
		return fmt.Errorf("event ID is required")
	}

	if event.Type == "" {
		return fmt.Errorf("event type is required")
	}

	event.CreatedAt = time.Now().UTC()

	if err := s.repo.Save(ctx, event); err != nil {
		return fmt.Errorf("saving event: %w", err)
	}

	s.logger.Info("event processed",
		"event_id", event.ID,
		"event_type", event.Type,
	)
	return nil
}
