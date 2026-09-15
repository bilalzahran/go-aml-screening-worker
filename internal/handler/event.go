package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/service"
)

type EventHandler struct {
	service *service.EventService
	logger  *slog.Logger
}

func NewEventHandler(svc *service.EventService, logger *slog.Logger) *EventHandler {
	return &EventHandler{
		service: svc,
		logger:  logger.With("component", "event_handler"),
	}
}

func (h *EventHandler) Handle(ctx context.Context, body []byte) error {
	var event model.Event
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshaling event: %w", err)
	}

	if err := h.service.Process(ctx, &event); err != nil {
		return fmt.Errorf("processing event: %w", err)
	}

	return nil
}
