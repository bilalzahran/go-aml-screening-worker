package repository

import (
	"cdaq-event-worker/internal/model"
	"context"
)

type TenantRepository interface {
	FindByPubId(ctx context.Context, pubId string) (*model.Tenant, error)
}
