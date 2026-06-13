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

func hashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h), err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func validateCredentialsUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 2 {
		return "", fmt.Errorf("username too short")
	}
	runes := []rune(username)
	if len(runes) > 32 {
		username = string(runes[:32])
	}
	return username, nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password too short")
	}
	if len(password) > 128 {
		return fmt.Errorf("password too long")
	}
	return nil
}

// SignInWithPassword authenticates by username and password.
func (s *Service) SignInWithPassword(ctx context.Context, username, password string) (string, *domain.UserMe, error) {
	username, err := validateCredentialsUsername(username)
	if err != nil {
		return "", nil, err
	}
	if err := validatePassword(password); err != nil {
		return "", nil, err
	}

	u, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		return "", nil, ErrUserNotExists
	}
	if !checkPassword(u.HashedPassword, password) {
		return "", nil, ErrWrongPassword
	}

	token, err := s.jwt.GenerateToken(u.ID, u.Username, true)
	if err != nil {
		return "", nil, err
	}
	me := s.UserMeFromUser(u)
	return token, &me, nil
}

// SignUpWithPassword registers a new account with username and password.
func (s *Service) SignUpWithPassword(ctx context.Context, username, password string) (string, *domain.UserMe, error) {
	if !s.signupAllowed {
		return "", nil, ErrSignupClosed
	}

	username, err := validateCredentialsUsername(username)
	if err != nil {
		return "", nil, err
	}
	if err := validatePassword(password); err != nil {
		return "", nil, err
	}

	taken, err := s.storage.IsUsernameTaken(ctx, username, uuid.Nil)
	if err != nil {
		return "", nil, err
	}
	if taken {
		return "", nil, ErrUsernameTaken
	}

	pwHash, err := hashPassword(password)
	if err != nil {
		return "", nil, err
	}

	user := &domain.User{
		ID:             uuid.New(),
		Username:       username,
		HashedPassword: pwHash,
	}
	if err := s.storage.CreateUser(ctx, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", nil, ErrUsernameTaken
		}
		return "", nil, err
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Username, true)
	if err != nil {
		return "", nil, err
	}
	me := s.UserMeFromUser(user)
	return token, &me, nil
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
