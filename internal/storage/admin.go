package storage

import (
	"context"
	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
)

func (s *Storage) CreateContest(ctx context.Context, c *domain.Contest) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO contests (id, type, year, starts, ends)
		VALUES ($1, $2, $3, $4, $5)
	`, c.ID, c.Type, c.Year, c.Starts, c.Ends)
	return err
}

func (s *Storage) UpdateContest(ctx context.Context, c *domain.Contest) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE contests
		SET type = $2, year = $3, starts = $4, ends = $5
		WHERE id = $1
	`, c.ID, c.Type, c.Year, c.Starts, c.Ends)
	return err
}

func (s *Storage) CreatePerformance(ctx context.Context, p *domain.Performance, contestID, countryID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	nextNumber := p.Number
	if nextNumber <= 0 {
		err := s.pool.QueryRow(ctx, `
			SELECT COALESCE(MAX(number), 0) + 1 FROM performance WHERE contest_id = $1
		`, contestID).Scan(&nextNumber)
		if err != nil {
			return uuid.Nil, err
		}
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO performance (id, contest_id, country_id, number, youtube_link, artist, song, qualified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, false)
	`, id, contestID, countryID, nextNumber, p.YoutubeLink, p.Artist, p.Song)
	return id, err
}

func (s *Storage) UpdateContestPerformancePlaces(ctx context.Context, contestID uuid.UUID, ranked []uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE performance SET place = NULL WHERE contest_id = $1`, contestID)
	if err != nil {
		return err
	}

	for i, perfID := range ranked {
		place := i + 1
		_, err = tx.Exec(ctx, `
			UPDATE performance SET place = $2 WHERE id = $1 AND contest_id = $3
		`, perfID, place, contestID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Storage) UpdatePerformanceFields(ctx context.Context, performanceID uuid.UUID, qualified bool, link string, place *int) error {
	if place == nil {
		_, err := s.pool.Exec(ctx, `
			UPDATE performance
			SET qualified = $1, youtube_link = $2, place = NULL
			WHERE id = $3
		`, qualified, link, performanceID)
		return err
	}
	_, err := s.pool.Exec(ctx, `
		UPDATE performance
		SET qualified = $1, youtube_link = $2, place = $4
		WHERE id = $3
	`, qualified, link, performanceID, *place)
	return err
}
