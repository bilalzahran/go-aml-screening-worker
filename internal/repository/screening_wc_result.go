package repository

import (
	"context"

	"cdaq-event-worker/internal/model"
)

type ScreeningWcResultRepository interface {
	Save(ctx context.Context, tenantId string, result *model.ScreeningWcResult) error
}
