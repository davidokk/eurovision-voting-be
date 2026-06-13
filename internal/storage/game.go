package storage

import (
	"context"
	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
)

func (s *Storage) ListGameCatalog(ctx context.Context) ([]domain.GameCatalogItem, error) {
	query := `
		SELECT
			p.id::text,
			p.artist,
			p.song,
			c.name_ru,
			c.flag_emogi,
			co.year,
			co.type,
			p.youtube_link
		FROM performance p
		JOIN countries c ON c.id = p.country_id
		JOIN contests co ON co.id = p.contest_id
		WHERE p.youtube_link IS NOT NULL AND p.youtube_link <> ''
		ORDER BY co.year DESC, co.type, p.number ASC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.GameCatalogItem
	for rows.Next() {
		var item domain.GameCatalogItem
		if err := rows.Scan(
			&item.PerformanceID,
			&item.Artist,
			&item.Song,
			&item.CountryName,
			&item.FlagEmoji,
			&item.Year,
			&item.ContestType,
			&item.YoutubeLink,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Storage) GetGameCatalogItems(ctx context.Context, ids []string) ([]domain.GameCatalogItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := `
		SELECT
			p.id::text,
			p.artist,
			p.song,
			c.name_ru,
			c.flag_emogi,
			co.year,
			co.type,
			p.youtube_link
		FROM performance p
		JOIN countries c ON c.id = p.country_id
		JOIN contests co ON co.id = p.contest_id
		WHERE p.id::text = ANY($1)
	`
	rows, err := s.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byID := make(map[string]domain.GameCatalogItem)
	for rows.Next() {
		var item domain.GameCatalogItem
		if err := rows.Scan(
			&item.PerformanceID,
			&item.Artist,
			&item.Song,
			&item.CountryName,
			&item.FlagEmoji,
			&item.Year,
			&item.ContestType,
			&item.YoutubeLink,
		); err != nil {
			return nil, err
		}
		byID[item.PerformanceID] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	ordered := make([]domain.GameCatalogItem, 0, len(ids))
	for _, id := range ids {
		if item, ok := byID[id]; ok {
			ordered = append(ordered, item)
		}
	}
	return ordered, nil
}

func (s *Storage) GetScoresForPerformance(ctx context.Context, performanceID uuid.UUID) ([]domain.GameContestScore, error) {
	query := `
		SELECT u.username, sc.score, sc.comment
		FROM scores sc
		JOIN users u ON u.id = sc.user_id
		WHERE sc.performance_id = $1
		ORDER BY sc.score DESC, u.username ASC
	`
	rows, err := s.pool.Query(ctx, query, performanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.GameContestScore
	for rows.Next() {
		var item domain.GameContestScore
		if err := rows.Scan(&item.Username, &item.Score, &item.Comment); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
