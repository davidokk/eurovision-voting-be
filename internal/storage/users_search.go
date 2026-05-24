package storage

import (
	"context"
	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) ListUsers(ctx context.Context, excludeID *uuid.UUID, limit int) ([]domain.UserPublic, error) {
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	var (
		rows pgx.Rows
		err  error
	)
	if excludeID != nil {
		rows, err = s.pool.Query(ctx, `
			SELECT id::text, username, avatar_url
			FROM users
			WHERE id != $1
			ORDER BY username
			LIMIT $2
		`, *excludeID, limit)
	} else {
		rows, err = s.pool.Query(ctx, `
			SELECT id::text, username, avatar_url
			FROM users
			ORDER BY username
			LIMIT $1
		`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.UserPublic
	for rows.Next() {
		var u domain.UserPublic
		if err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, rows.Err()
}
