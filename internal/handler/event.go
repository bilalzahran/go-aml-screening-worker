package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/service"
)

type EventTypeHandler interface {
	Handle(ctx context.Context, event *model.Event) error
}

type EventHandler struct {
	handlers     map[string]EventTypeHandler
	eventService *service.EventService
	logger       *slog.Logger
}

func NewEventHandler(handlers map[string]EventTypeHandler, eventService *service.EventService, logger *slog.Logger) *EventHandler {
	return &EventHandler{
		handlers:     handlers,
		eventService: eventService,
		logger:       logger.With("component", "event_handler"),
	}
}

func (h *EventHandler) Handle(ctx context.Context, body []byte) error {
	var event model.Event
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshaling event: %w", err)
	}

	handler, ok := h.handlers[event.Type]
	if !ok {
		return fmt.Errorf("no handler for event type: %s", event.Type)
	}

	// Process event - save to db
	h.eventService.Process(ctx, &event)

	return handler.Handle(ctx, &event)
}
