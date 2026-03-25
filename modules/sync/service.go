package sync

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type service struct {
	repo *repository
}

func NewService(database db.Querier) Service {
	return &service{repo: NewRepository(database)}
}

func (s *service) Pull(ctx context.Context, lastSync time.Time) (*SyncBatch, error) {
	// In a real scenario, enterpriseID would come from context
	// For now, I'll use a placeholder or assume it's set in ctx
	enterpriseID := int64(1) 

	productsList, err := s.repo.GetProductUpdates(ctx, lastSync, enterpriseID)
	if err != nil {
		return nil, err
	}

	thirdPartiesList, err := s.repo.GetThirdPartyUpdates(ctx, lastSync, enterpriseID)
	if err != nil {
		return nil, err
	}

	return &SyncBatch{
		Products:     productsList,
		ThirdParties: thirdPartiesList,
	}, nil
}

func (s *service) Push(ctx context.Context, batch *SyncBatch) (*SyncStats, error) {
	start := time.Now()
	stats := &SyncStats{}

	for _, p := range batch.Products {
		if err := s.repo.UpsertProduct(ctx, p); err != nil {
			stats.FailedCount++
			continue
		}
		stats.PushedCount++
	}

	for _, tp := range batch.ThirdParties {
		if err := s.repo.UpsertThirdParty(ctx, tp); err != nil {
			stats.FailedCount++
			continue
		}
		stats.PushedCount++
	}

	stats.Duration = time.Since(start)
	return stats, nil
}
