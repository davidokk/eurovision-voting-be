package service

import (
	"context"
	"eurovision-voting/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func (s *Service) RatePerformance(ctx context.Context, userID, performanceID uuid.UUID, score int, comment, gif string) error {
	c, err := s.storage.GetContestByPerformance(ctx, performanceID)
	if err != nil {
		return err
	}
	if time.Now().Before(c.Starts) || time.Now().After(c.Ends) {
		return ErrContestClosed
	}
	if err := s.storage.RatePerformance(ctx, userID, performanceID, score, comment, gif); err != nil {
		return err
	}
	if err := s.storage.InsertMessage(ctx, &domain.Message{
		Type:          "system",
		UserID:        userID,
		PerformanceID: performanceID,
		CreatedAt:     time.Now(),
		ContestID:     c.ID,
	}); err != nil {
		log.Error().Err(err).Msg("cannot insert score message")
	}
	country, err := s.storage.GetCountryByPerformance(ctx, performanceID)
	if err != nil {
		log.Error().Err(err).Msg("cannot get country")
		return nil 
	}
	u, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("cannot get user")
		return nil 
	}
	s.broadcastMessage(&domain.Message{
		Type: "system",
		UserID: userID,
		Username: u.Username,
		ContestID: c.ID,
		PerformanceID: performanceID,
		CreatedAt: time.Now(),
		Score: &score,
		Gif: &gif,
		Comment: &comment,
		Country: &country.NameRU,
		CountryFlag: &country.FlagEmoji,
	})
	return nil
}

func (s *Service) UpdatePerformance(ctx context.Context, id uuid.UUID, qualified bool, link string) error {
	return s.storage.UpdatePerformance(ctx, id, qualified, link)
}

func (s *Service) GetScoresFiltered(ctx context.Context, f domain.Filters) ([]domain.ScoreFiltered, error) {
	return s.storage.GetScoresFiltered(ctx, f)
}

func (s *Service) GetCountries(ctx context.Context) ([]domain.Country, error) {
	return s.storage.GetCountries(ctx)
}
