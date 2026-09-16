package usecase

import (
	"cdaq-event-worker/internal/model"
	"context"
	"log/slog"
)

type JobProcessor interface {
	CreateJob(ctx context.Context, event *model.Event) error
}

type OgsUseCase struct {
	logger       *slog.Logger
	JobProcessor JobProcessor
}

func NewOgsUseCase(logger *slog.Logger, jobProcessor JobProcessor) *OgsUseCase {
	return &OgsUseCase{
		logger: logger,
	}
}

func (o *OgsUseCase) Handle(ctx context.Context, event *model.Event) error {
	return nil
}
