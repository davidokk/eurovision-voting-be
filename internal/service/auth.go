package service

import (
	"context"
	"eurovision-voting/internal/domain"
)

// SignUp and SignIn are replaced by email verification flows in auth_email.go.
// Kept for compatibility with middleware that loads users by username from JWT.

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.storage.GetUserByUsername(ctx, username)
}
