package storage

import (
	"context"
	"time"

	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
)

func (s *Storage) AddPerformanceFavorite(ctx context.Context, userID, performanceID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO performance_favorites (user_id, performance_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, performance_id) DO NOTHING
	`, userID, performanceID)
	return err
}

func (s *Storage) RemovePerformanceFavorite(ctx context.Context, userID, performanceID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		DELETE FROM performance_favorites
		WHERE user_id = $1 AND performance_id = $2
	`, userID, performanceID)
	return err
}

func (s *Storage) ListFavoritePerformanceIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT performance_id::text
		FROM performance_favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (s *Storage) ListFavoritePerformances(ctx context.Context, userID uuid.UUID) ([]domain.FavoritePerformance, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT
			p.id::text,
			c.name_ru,
			c.flag_emogi,
			co.year,
			co.type,
			p.song,
			p.artist,
			p.youtube_link,
			p.qualified,
			p.place,
			f.created_at
		FROM performance_favorites f
		JOIN performance p ON p.id = f.performance_id
		JOIN countries c ON c.id = p.country_id
		JOIN contests co ON co.id = p.contest_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.FavoritePerformance
	for rows.Next() {
		var item domain.FavoritePerformance
		var createdAt time.Time
		if err := rows.Scan(
			&item.PerformanceID,
			&item.CountryName,
			&item.FlagEmoji,
			&item.ContestYear,
			&item.ContestType,
			&item.Song,
			&item.Artist,
			&item.YoutubeLink,
			&item.Qualified,
			&item.Place,
			&createdAt,
		); err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt
		out = append(out, item)
	}
	return out, rows.Err()
}
