package service

import (
	"context"
	"errors"
	"eurovision-voting/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SignUp(ctx context.Context, username, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("encrypt password: %w", err)
	}

	user := &domain.User{
		ID:             uuid.New(),
		Username:       username,
		HashedPassword: string(passwordHash),
	}

	if err := s.storage.CreateUser(ctx, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", ErrUserAlreadyExists
			}
		}
		return "", err
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (s *Service) SignIn(ctx context.Context, username, password string) (string, error) {
	u, err := s.storage.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotExists
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)); err != nil {
		return "", ErrWrongPassword
	}

	token, err := s.jwt.GenerateToken(u.ID, u.Username)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.storage.GetUserByUsername(ctx, username)
}