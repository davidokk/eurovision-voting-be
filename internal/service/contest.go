package service

import (
	"context"
	"eurovision-voting/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Service) GetContestView(ctx context.Context, contestID uuid.UUID) (*domain.ContestView, error) {
	return s.storage.GetContestView(ctx, contestID)
}

func (s *Service) GetContestsByYears(ctx context.Context) (map[int][]domain.Contest, error) {
	cts, err := s.storage.GetContestList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get contest list: %w", err)
	}

	res := map[int][]domain.Contest{}
	for _, c := range cts {
		if _, ok := res[c.Year]; !ok {
			res[c.Year] = []domain.Contest{}
		}
		res[c.Year] = append(res[c.Year], c)
	}
	return res, nil
}

func (s *Service) RatePerformance(ctx context.Context, userID, performanceID uuid.UUID, score int, comment string) error {
	c, err := s.storage.GetContestByPerformance(ctx, performanceID)
	if err != nil {
		return err
	}
	if time.Now().Before(c.Starts) || time.Now().After(c.Ends) {
		return ErrContestClosed
	}
	return s.storage.RatePerformance(ctx, userID, performanceID, score, comment)
}
