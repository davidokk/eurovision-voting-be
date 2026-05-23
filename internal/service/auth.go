package service

import (
	"context"
	"crypto/rand"
	"errors"
	"eurovision-voting/internal/domain"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.storage.GetUserByUsername(ctx, username)
}

func (s *Service) IsTelegramLinked(u *domain.User) bool {
	return u.TelegramID != nil && *u.TelegramID > 0
}

func (s *Service) UserMeFromUser(u *domain.User) domain.UserMe {
	tgUser := ""
	if u.TelegramUsername != nil {
		tgUser = *u.TelegramUsername
	}
	return domain.UserMe{
		UserPublic: domain.UserPublic{
			ID:        u.ID.String(),
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		},
		TelegramLinked:   s.IsTelegramLinked(u),
		TelegramUsername: tgUser,
	}
}

func (s *Service) issueToken(u *domain.User) (string, error) {
	return s.jwt.GenerateToken(u.ID, u.Username, s.IsTelegramLinked(u))
}

func generateSixDigitCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func hashCode(code string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	return string(b), err
}

func checkCode(hash, code string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(code)) == nil
}

func (s *Service) ChangeUsername(ctx context.Context, userID uuid.UUID, username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 2 {
		return "", fmt.Errorf("username too short")
	}

	taken, err := s.storage.IsUsernameTaken(ctx, username, userID)
	if err != nil {
		return "", err
	}
	if taken {
		return "", ErrUsernameTaken
	}

	if err := s.storage.UpdateUsername(ctx, userID, username); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", ErrUsernameTaken
		}
		return "", err
	}

	u, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}
	return s.issueToken(u)
}
