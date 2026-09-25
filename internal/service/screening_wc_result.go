package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/repository"
)

type ScreeningWcResultService struct {
	repo   repository.ScreeningWcResultRepository
	logger *slog.Logger
}

func NewScreeningWcResultService(repo repository.ScreeningWcResultRepository, logger *slog.Logger) *ScreeningWcResultService {
	return &ScreeningWcResultService{
		repo:   repo,
		logger: logger.With("component", "screening_wc_result_service"),
	}
}

func (s *ScreeningWcResultService) Save(ctx context.Context, result *model.ScreeningWcResult, tenantId string) error {
	// Generate new ID if empty
	if result.BaseEntity.ID.IsZero() {
		result.BaseEntity.ID = bson.NewObjectID()
	}

	// Set timestamps
	now := time.Now().UTC()
	result.BaseEntity.CreatedAt = now
	result.BaseEntity.UpdatedAt = now

	if err := s.repo.Save(ctx, tenantId, result); err != nil {
		return fmt.Errorf("saving screening_wc_result: %w", err)
	}

	s.logger.Info("screening_wc_result saved",
		"id", result.BaseEntity.ID,
	)
	return nil
}
