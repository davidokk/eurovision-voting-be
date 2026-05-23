package service

import "errors"

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotExists       = errors.New("user not exists")
	ErrWrongPassword       = errors.New("wrong password")
	ErrContestClosed       = errors.New("Контесты закрыт, нельзя оценить")
	ErrInvalidCode         = errors.New("invalid or expired code")
	ErrEmailTaken          = errors.New("email already registered")
	ErrUsernameTaken       = errors.New("username already taken")
	ErrEmailNotConfigured  = errors.New("email service not configured")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrEmailRateLimit      = errors.New("email rate limit exceeded")
)
