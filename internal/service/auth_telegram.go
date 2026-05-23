package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"eurovision-voting/internal/domain"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

const (
	telegramPurposeSignin = "signin"

	telegramSessionTTL      = 15 * time.Minute
	telegramRateLimitMax      = 5
	telegramRateLimitWindow = time.Hour
	telegramSyntheticEmailFmt = "tg_%d@telegram.local"
)

var telegramUsernameSanitize = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

type TelegramStartResult struct {
	LinkToken   string `json:"link_token"`
	BotURL      string `json:"bot_url"`
	BotUsername string `json:"bot_username,omitempty"`
}

type TelegramSessionStatus struct {
	Status            string `json:"status"`
	TelegramConnected bool   `json:"telegram_connected"`
	CodeSent          bool   `json:"code_sent"`
}

func (s *Service) TelegramConfigured() bool {
	return s.telegramBotUsername != ""
}

func (s *Service) TelegramBotURL(linkToken string) string {
	if s.telegramBotUsername == "" {
		return ""
	}
	return fmt.Sprintf("https://t.me/%s?start=%s", s.telegramBotUsername, linkToken)
}

func newLinkToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func randomPasswordHash() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	h, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	return string(h), err
}

func sanitizeSiteUsername(from string) string {
	from = strings.TrimSpace(from)
	if from == "" {
		return ""
	}
	from = strings.TrimPrefix(from, "@")
	from = telegramUsernameSanitize.ReplaceAllString(from, "")
	from = strings.Trim(from, "_")
	if len(from) < 2 {
		return ""
	}
	runes := []rune(from)
	if len(runes) > 32 {
		from = string(runes[:32])
	}
	return from
}

func (s *Service) allocateUsername(ctx context.Context, tgUsername string, telegramID int64) (string, error) {
	candidates := []string{}
	if base := sanitizeSiteUsername(tgUsername); base != "" {
		candidates = append(candidates, base)
	}
	candidates = append(candidates, fmt.Sprintf("tg%d", telegramID))

	for _, name := range candidates {
		taken, err := s.storage.IsUsernameTaken(ctx, name, uuid.Nil)
		if err != nil {
			return "", err
		}
		if !taken {
			return name, nil
		}
	}

	base := candidates[0]
	for i := 0; i < 20; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(9000))
		if err != nil {
			return "", err
		}
		name := fmt.Sprintf("%s_%d", base, n.Int64()+1000)
		taken, err := s.storage.IsUsernameTaken(ctx, name, uuid.Nil)
		if err != nil {
			return "", err
		}
		if !taken {
			return name, nil
		}
	}
	return "", ErrUsernameTaken
}

func (s *Service) checkTelegramRateLimit(ctx context.Context, telegramID int64) error {
	since := time.Now().Add(-telegramRateLimitWindow)
	n, err := s.storage.CountTelegramCodesSince(ctx, telegramID, since)
	if err != nil {
		return err
	}
	if n >= telegramRateLimitMax {
		return ErrTelegramRateLimit
	}
	return nil
}

func (s *Service) createTelegramSession(ctx context.Context) (*TelegramStartResult, error) {
	if !s.TelegramConfigured() {
		return nil, ErrTelegramNotConfigured
	}

	linkToken, err := newLinkToken()
	if err != nil {
		return nil, err
	}

	sess := &domain.TelegramAuthSession{
		ID:        uuid.New(),
		LinkToken: linkToken,
		Purpose:   telegramPurposeSignin,
		ExpiresAt: time.Now().Add(telegramSessionTTL),
	}
	if err := s.storage.CreateTelegramAuthSession(ctx, sess); err != nil {
		return nil, err
	}

	return &TelegramStartResult{
		LinkToken:   linkToken,
		BotURL:      s.TelegramBotURL(linkToken),
		BotUsername: s.telegramBotUsername,
	}, nil
}

func (s *Service) StartTelegramSignin(ctx context.Context) (*TelegramStartResult, error) {
	return s.createTelegramSession(ctx)
}

func (s *Service) GetTelegramSessionStatus(ctx context.Context, linkToken string) (*TelegramSessionStatus, error) {
	sess, err := s.storage.GetTelegramAuthSessionByToken(ctx, linkToken)
	if err != nil {
		return nil, ErrTelegramSessionInvalid
	}
	st := &TelegramSessionStatus{Status: "pending"}
	if sess.TelegramID != nil {
		st.TelegramConnected = true
		st.Status = "awaiting_code"
	}
	if sess.CodeSentAt != nil {
		st.CodeSent = true
	}
	return st, nil
}

// BotDeliverCode attaches Telegram user to session and stores hashed code; returns plaintext code for the bot message.
func (s *Service) BotDeliverCode(ctx context.Context, linkToken string, telegramID, chatID int64, telegramUsername string) (string, error) {
	sess, err := s.storage.GetTelegramAuthSessionByToken(ctx, linkToken)
	if err != nil {
		return "", ErrTelegramSessionInvalid
	}
	if sess.Purpose != telegramPurposeSignin && sess.Purpose != "signup" && sess.Purpose != "link" {
		return "", ErrTelegramSessionInvalid
	}

	_, userErr := s.storage.GetUserByTelegramID(ctx, telegramID)
	if errors.Is(userErr, pgx.ErrNoRows) && !s.signupAllowed {
		return "", ErrSignupClosed
	}
	if userErr != nil && !errors.Is(userErr, pgx.ErrNoRows) {
		return "", userErr
	}

	if err := s.checkTelegramRateLimit(ctx, telegramID); err != nil {
		return "", err
	}

	_ = s.storage.InvalidateTelegramSessions(ctx, sess.Purpose, &telegramID, nil)

	code, err := generateSixDigitCode()
	if err != nil {
		return "", err
	}
	codeHash, err := hashCode(code)
	if err != nil {
		return "", err
	}

	if err := s.storage.AttachTelegramToSession(ctx, sess.ID, telegramID, chatID, telegramUsername, codeHash); err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) findTelegramSessionCode(sess *domain.TelegramAuthSession, code string) bool {
	if sess.CodeHash == nil {
		return false
	}
	return checkCode(*sess.CodeHash, code)
}

func (s *Service) createUserFromTelegramSession(ctx context.Context, sess *domain.TelegramAuthSession) (*domain.User, error) {
	if sess.TelegramID == nil {
		return nil, ErrTelegramNotConnected
	}

	tgID := *sess.TelegramID
	tgName := ""
	if sess.TelegramUsername != nil {
		tgName = *sess.TelegramUsername
	}

	username, err := s.allocateUsername(ctx, tgName, tgID)
	if err != nil {
		return nil, err
	}

	pwHash, err := randomPasswordHash()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	email := fmt.Sprintf(telegramSyntheticEmailFmt, tgID)
	user := &domain.User{
		ID:               uuid.New(),
		Username:         username,
		Email:            email,
		EmailVerifiedAt:  &now,
		TelegramID:       &tgID,
		TelegramUsername: sess.TelegramUsername,
		TelegramLinkedAt: &now,
		HashedPassword:   pwHash,
	}

	if err := s.storage.CreateUser(ctx, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "telegram") {
				return nil, ErrTelegramAlreadyLinked
			}
			return nil, ErrUsernameTaken
		}
		return nil, err
	}
	return user, nil
}

// ConfirmTelegramSignin logs in an existing user or creates a new account when signup is allowed.
func (s *Service) ConfirmTelegramSignin(ctx context.Context, linkToken, code string, signupAllowed bool) (string, *domain.UserMe, error) {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return "", nil, ErrInvalidCode
	}

	sess, err := s.storage.GetTelegramAuthSessionByToken(ctx, linkToken)
	if err != nil {
		return "", nil, ErrTelegramSessionInvalid
	}
	if sess.TelegramID == nil {
		return "", nil, ErrTelegramNotConnected
	}
	if sess.Purpose != telegramPurposeSignin && sess.Purpose != "signup" && sess.Purpose != "link" {
		return "", nil, ErrTelegramSessionInvalid
	}
	if !s.findTelegramSessionCode(sess, code) {
		return "", nil, ErrInvalidCode
	}

	var u *domain.User
	existing, err := s.storage.GetUserByTelegramID(ctx, *sess.TelegramID)
	if err == nil {
		u = existing
	} else if errors.Is(err, pgx.ErrNoRows) {
		if !signupAllowed {
			return "", nil, ErrSignupClosed
		}
		u, err = s.createUserFromTelegramSession(ctx, sess)
		if errors.Is(err, ErrTelegramAlreadyLinked) {
			u, err = s.storage.GetUserByTelegramID(ctx, *sess.TelegramID)
		}
		if err != nil {
			return "", nil, err
		}
	} else {
		return "", nil, err
	}

	_ = s.storage.MarkTelegramSessionUsed(ctx, sess.ID)

	token, err := s.issueToken(u)
	if err != nil {
		return "", nil, err
	}
	me := s.UserMeFromUser(u)
	return token, &me, nil
}
