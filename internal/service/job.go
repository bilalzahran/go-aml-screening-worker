package service

import (
	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
	"context"
	"log/slog"
)

type JobService struct {
	jobRepo      repository.JobRepository
	childJobRepo repository.ChildJobRepository
	logger       *slog.Logger
}

func NewJobService(jobRepo repository.JobRepository, childJobRepo repository.ChildJobRepository, logger *slog.Logger) *JobService {
	return &JobService{
		jobRepo:      jobRepo,
		childJobRepo: childJobRepo,
		logger:       logger,
	}
}

func (j *JobService) CreateJob(ctx context.Context, event *model.Event) error {
	return nil
}
