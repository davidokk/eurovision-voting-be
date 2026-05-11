package service

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotExists = errors.New("user not exists")
	ErrWrongPassword = errors.New("wrong password")
	ErrContestClosed = errors.New("Контесты закрыт, нельзя оценить")
)
