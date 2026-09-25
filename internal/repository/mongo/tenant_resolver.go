package mongo

import (
	"cdaq-event-worker/internal/repository"
	"context"
	"log/slog"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/sync/singleflight"
)

type TenantResolver struct {
	tenantRepo repository.TenantRepository
	client     *mongo.Client
	dbs        map[string]*mongo.Database
	mu         sync.RWMutex
	sf         singleflight.Group
	logger     *slog.Logger
}

func NewTenantResolver(tenantRepo repository.TenantRepository, mongoClient *mongo.Client, logger *slog.Logger) *TenantResolver {
	return &TenantResolver{
		client:     mongoClient,
		tenantRepo: tenantRepo,
		logger:     logger,
		dbs:        make(map[string]*mongo.Database),
	}
}

func (r *TenantResolver) Database(ctx context.Context, tenantPubId string) (*mongo.Database, error) {
	r.mu.RLock()
	db, ok := r.dbs[tenantPubId]
	r.mu.RUnlock()
	if ok {
		return db, nil
	}

	// do single flight
	result, err, _ := r.sf.Do(tenantPubId, func() (any, error) {
		tenant, err := r.tenantRepo.FindByPubId(ctx, tenantPubId)
		if err != nil {
			return nil, err
		}
		return r.client.Database(tenant.DbName), nil
	})

	if err != nil {
		return nil, err
	}

	resultDb := result.(*mongo.Database)

	r.mu.Lock()
	r.dbs[tenantPubId] = resultDb
	r.mu.Unlock()

	return resultDb, nil
}

func (r *TenantResolver) Close() error {
	return nil

}
