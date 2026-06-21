package service

import (
	"context"
	"errors"

	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
)

var ErrFavoriteNotFound = errors.New("favorite not found")

func (s *Service) ListFavoritePerformanceIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.storage.ListFavoritePerformanceIDs(ctx, userID)
}

func (s *Service) ListFavoritePerformances(ctx context.Context, userID uuid.UUID) ([]domain.FavoritePerformance, error) {
	return s.storage.ListFavoritePerformances(ctx, userID)
}

func (s *Service) AddPerformanceFavorite(ctx context.Context, userID, performanceID uuid.UUID) error {
	return s.storage.AddPerformanceFavorite(ctx, userID, performanceID)
}

func (s *Service) RemovePerformanceFavorite(ctx context.Context, userID, performanceID uuid.UUID) error {
	return s.storage.RemovePerformanceFavorite(ctx, userID, performanceID)
}
