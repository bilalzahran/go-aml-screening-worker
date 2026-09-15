package repository

import (
	"context"

	"cdaq-event-worker/internal/model"
)

type EventRepository interface {
	Save(ctx context.Context, event *model.Event) error
}
