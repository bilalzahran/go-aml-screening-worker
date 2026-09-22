package usecase

import (
	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/service"
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

const MAX_CONCURRENCY = 50

type OgsUseCase struct {
	logger                   *slog.Logger
	jobService               *service.JobService
	screeningWcResultService *service.ScreeningWcResultService
}

func NewOgsUseCase(logger *slog.Logger, jobService *service.JobService, screeningWcResultService *service.ScreeningWcResultService) *OgsUseCase {
	return &OgsUseCase{
		logger:                   logger,
		jobService:               jobService,
		screeningWcResultService: screeningWcResultService,
	}
}

func (o *OgsUseCase) Handle(ctx context.Context, event *model.Event) error {
	results := make([]model.ScreeningWcResult, len(event.WorldCheck))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(MAX_CONCURRENCY)

	for i, child := range event.WorldCheck {
		childEvent, err := model.MapToWorldCheckHits(child)
		if err != nil {
			return err
		}

		g.Go(func() error {
			res, _ := o.processOneChild(ctx, childEvent)
			results[i] = *res
			return nil
		})
	}

	// TODO: Final save to screening wc + screening wc ogs

	return nil
}

func (o *OgsUseCase) processOneChild(ctx context.Context, event *model.WorldCheckHits) (*model.ScreeningWcResult, error) {
	screeningWcResult := model.NewScreeningWcResultFromWorldCheckHits(event)

	if err := o.screeningWcResultService.Save(ctx, screeningWcResult); err != nil {
		return nil, err
	}

	// TODO: Call Typesafe API / Jev

	return screeningWcResult, nil
}
