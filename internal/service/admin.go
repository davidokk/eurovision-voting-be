package service

import (
	"context"
	"eurovision-voting/internal/domain"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Service) AdminCreateContest(ctx context.Context, year int, contestType string, starts, ends time.Time) (*domain.Contest, error) {
	contestType = strings.TrimSpace(contestType)
	switch contestType {
	case "final", "first-semifinal", "second-semifinal":
	default:
		return nil, fmt.Errorf("invalid contest type")
	}
	if ends.Before(starts) {
		return nil, fmt.Errorf("ends must be after starts")
	}
	c := &domain.Contest{
		ID:     uuid.New(),
		Type:   contestType,
		Year:   year,
		Starts: starts,
		Ends:   ends,
	}
	if err := s.storage.CreateContest(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) AdminUpdateContest(ctx context.Context, id uuid.UUID, year int, contestType string, starts, ends time.Time) (*domain.Contest, error) {
	contestType = strings.TrimSpace(contestType)
	switch contestType {
	case "final", "first-semifinal", "second-semifinal":
	default:
		return nil, fmt.Errorf("invalid contest type")
	}
	if ends.Before(starts) {
		return nil, fmt.Errorf("ends must be after starts")
	}
	c := &domain.Contest{
		ID:     id,
		Type:   contestType,
		Year:   year,
		Starts: starts,
		Ends:   ends,
	}
	if err := s.storage.UpdateContest(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) AdminCreatePerformance(
	ctx context.Context,
	contestID uuid.UUID,
	countryID uuid.UUID,
	number int,
	artist, song, youtubeLink string,
) (uuid.UUID, error) {
	p := &domain.Performance{
		Number:      number,
		Artist:      strings.TrimSpace(artist),
		Song:        strings.TrimSpace(song),
		YoutubeLink: strings.TrimSpace(youtubeLink),
	}
	if p.Artist == "" || p.Song == "" {
		return uuid.Nil, fmt.Errorf("artist and song required")
	}
	if p.YoutubeLink == "" {
		p.YoutubeLink = "https://youtube.com/"
	}
	return s.storage.CreatePerformance(ctx, p, contestID, countryID)
}

func (s *Service) AdminUpdatePlaces(ctx context.Context, contestID uuid.UUID, ranked []uuid.UUID) error {
	return s.storage.UpdateContestPerformancePlaces(ctx, contestID, ranked)
}

func (s *Service) UpdatePerformance(ctx context.Context, id uuid.UUID, qualified bool, link string, place *int) error {
	return s.storage.UpdatePerformanceFields(ctx, id, qualified, link, place)
}
