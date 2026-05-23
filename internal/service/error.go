package service

import "errors"

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserNotExists           = errors.New("user not exists")
	ErrWrongPassword           = errors.New("wrong password")
	ErrContestClosed           = errors.New("Контесты закрыт, нельзя оценить")
	ErrInvalidCode             = errors.New("invalid or expired code")
	ErrUsernameTaken           = errors.New("username already taken")
	ErrTelegramNotConfigured   = errors.New("telegram bot not configured")
	ErrTelegramRateLimit       = errors.New("telegram rate limit exceeded")
	ErrTelegramSessionInvalid  = errors.New("telegram session invalid")
	ErrTelegramNotConnected    = errors.New("telegram not connected to session")
	ErrTelegramAlreadyLinked   = errors.New("telegram already linked")
	ErrTelegramAccountNotFound = errors.New("telegram account not found")
	ErrSignupClosed            = errors.New("signup closed")
)
