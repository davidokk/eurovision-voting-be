package storage

import (
	"context"
	"eurovision-voting/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Storage) InvalidateAuthCodes(ctx context.Context, email, purpose string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE auth_codes SET used_at = now()
		WHERE lower(trim(email)) = lower(trim($1)) AND purpose = $2 AND used_at IS NULL
	`, email, purpose)
	return err
}

func (s *Storage) CreateAuthCode(ctx context.Context, c *domain.AuthCode) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO auth_codes (id, email, code_hash, purpose, user_id, username, password_hash, expires_at, created_at)
		VALUES ($1, lower(trim($2)), $3, $4, $5, $6, $7, $8, $9)
	`, c.ID, c.Email, c.CodeHash, c.Purpose, c.UserID, c.Username, c.PasswordHash, c.ExpiresAt, time.Now())
	return err
}

func (s *Storage) ListActiveAuthCodes(ctx context.Context, email, purpose string) ([]*domain.AuthCode, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, code_hash, purpose, user_id, username, password_hash, expires_at
		FROM auth_codes
		WHERE lower(trim(email)) = lower(trim($1))
		  AND purpose = $2
		  AND used_at IS NULL
		  AND expires_at > now()
		ORDER BY created_at DESC
		LIMIT 10
	`, email, purpose)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.AuthCode
	for rows.Next() {
		var c domain.AuthCode
		var userID *uuid.UUID
		var username, passwordHash *string
		if err := rows.Scan(&c.ID, &c.Email, &c.CodeHash, &c.Purpose, &userID, &username, &passwordHash, &c.ExpiresAt); err != nil {
			return nil, fmt.Errorf("scan auth code: %w", err)
		}
		c.UserID = userID
		c.Username = username
		c.PasswordHash = passwordHash
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (s *Storage) CountAuthCodesSince(ctx context.Context, email string, since time.Time) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM auth_codes
		WHERE lower(trim(email)) = lower(trim($1)) AND created_at >= $2
	`, email, since).Scan(&n)
	return n, err
}

func (s *Storage) MarkAuthCodeUsed(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `UPDATE auth_codes SET used_at = now() WHERE id = $1`, id)
	return err
}
