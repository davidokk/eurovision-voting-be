package service

import (
	"eurovision-voting/internal/jwt"
	"eurovision-voting/internal/storage"
)

type Service struct {
	storage *storage.Storage
	jwt     *jwt.Service
}

func New(storage *storage.Storage, jwt *jwt.Service) *Service {
	return &Service{
		storage: storage,
		jwt:     jwt,
	}
}
