package storage

import (
	"context"
	"eurovision-voting/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Storage) CreateTelegramAuthSession(ctx context.Context, sess *domain.TelegramAuthSession) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO telegram_auth_sessions (
			id, link_token, purpose, username, user_id, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, now())
	`, sess.ID, sess.LinkToken, sess.Purpose, sess.Username, sess.UserID, sess.ExpiresAt)
	return err
}

func (s *Storage) GetTelegramAuthSessionByToken(ctx context.Context, linkToken string) (*domain.TelegramAuthSession, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, link_token, purpose, username, user_id, telegram_id, telegram_chat_id,
		       telegram_username, code_hash, expires_at, code_sent_at
		FROM telegram_auth_sessions
		WHERE link_token = $1 AND used_at IS NULL AND expires_at > now()
		ORDER BY created_at DESC
		LIMIT 1
	`, linkToken)
	return scanTelegramSession(row)
}

func (s *Storage) AttachTelegramToSession(
	ctx context.Context,
	sessionID uuid.UUID,
	telegramID, chatID int64,
	telegramUsername string,
	codeHash string,
) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE telegram_auth_sessions
		SET telegram_id = $2,
		    telegram_chat_id = $3,
		    telegram_username = $4,
		    code_hash = $5,
		    code_sent_at = now()
		WHERE id = $1 AND used_at IS NULL
	`, sessionID, telegramID, chatID, nullIfEmpty(telegramUsername), codeHash)
	return err
}

func (s *Storage) MarkTelegramSessionUsed(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `UPDATE telegram_auth_sessions SET used_at = now() WHERE id = $1`, id)
	return err
}

func (s *Storage) InvalidateTelegramSessions(ctx context.Context, purpose string, telegramID *int64, userID *uuid.UUID) error {
	if telegramID != nil {
		_, err := s.pool.Exec(ctx, `
			UPDATE telegram_auth_sessions SET used_at = now()
			WHERE purpose = $1 AND telegram_id = $2 AND used_at IS NULL
		`, purpose, *telegramID)
		return err
	}
	if userID != nil {
		_, err := s.pool.Exec(ctx, `
			UPDATE telegram_auth_sessions SET used_at = now()
			WHERE purpose = $1 AND user_id = $2 AND used_at IS NULL
		`, purpose, *userID)
		return err
	}
	return nil
}

func (s *Storage) CountTelegramCodesSince(ctx context.Context, telegramID int64, since time.Time) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM telegram_auth_sessions
		WHERE telegram_id = $1 AND code_sent_at IS NOT NULL AND code_sent_at >= $2
	`, telegramID, since).Scan(&n)
	return n, err
}

func (s *Storage) GetUserByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `
		SELECT id, username, password, role, avatar_url,
		       coalesce(email, ''), email_verified_at,
		       telegram_id, telegram_username, telegram_linked_at
		FROM users WHERE telegram_id = $1
	`
	row := s.pool.QueryRow(ctx, query, telegramID)
	return scanUserRowFull(row)
}

func (s *Storage) LinkUserTelegram(ctx context.Context, userID uuid.UUID, telegramID int64, telegramUsername string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE users
		SET telegram_id = $2,
		    telegram_username = $3,
		    telegram_linked_at = now(),
		    email_verified_at = COALESCE(email_verified_at, now())
		WHERE id = $1
	`, userID, telegramID, nullIfEmpty(telegramUsername))
	return err
}

func (s *Storage) IsTelegramIDTaken(ctx context.Context, telegramID int64, excludeUserID *uuid.UUID) (bool, error) {
	var exists bool
	if excludeUserID != nil {
		err := s.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1 AND id <> $2)
		`, telegramID, *excludeUserID).Scan(&exists)
		return exists, err
	}
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)
	`, telegramID).Scan(&exists)
	return exists, err
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

type telegramSessionScanner interface {
	Scan(dest ...any) error
}

func scanTelegramSession(row telegramSessionScanner) (*domain.TelegramAuthSession, error) {
	var sess domain.TelegramAuthSession
	var username, tgUsername, codeHash *string
	var userID *uuid.UUID
	var telegramID, chatID *int64
	var codeSentAt *time.Time

	if err := row.Scan(
		&sess.ID, &sess.LinkToken, &sess.Purpose, &username, &userID,
		&telegramID, &chatID, &tgUsername, &codeHash, &sess.ExpiresAt, &codeSentAt,
	); err != nil {
		return nil, fmt.Errorf("scan telegram session: %w", err)
	}
	sess.Username = username
	sess.UserID = userID
	sess.TelegramID = telegramID
	sess.TelegramChatID = chatID
	sess.TelegramUsername = tgUsername
	sess.CodeHash = codeHash
	sess.CodeSentAt = codeSentAt
	return &sess, nil
}
