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
	/**
	TODO:
	1. Create worker (Fan-out)
	   a. Map the data
	   c. Fire the LLM -> do async, wait until the LLM gives back the result
	   d. Append the AI Recommendation to the hits
	   e. Save to db -> screening_wc_results collections
	   f. Return the _id of screening_wc_results
	2. Process the parent (Fan-in)
	   a. Update the screening_wc -> insert the results _id to the screening_wc_results
	   b. Update the screening_wc_ogs
	*/

	return nil
}

func (o *OgsUseCase) processOneChild(ctx context.Context, event *model.WorldCheckHits) error {
	return nil
}
