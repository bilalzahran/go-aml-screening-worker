package repository

import (
	"cdaq-event-worker/internal/model"
	"context"
)

type JobRepository interface {
	Save(ctx context.Context, job *model.Job) error
}

type ChildJobRepository interface {
	Save(ctx context.Context, childJob *model.ChildJob) error
}
