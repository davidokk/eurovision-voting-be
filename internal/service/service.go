package service

import (
	"eurovision-voting/internal/jwt"
	"eurovision-voting/internal/media"
	"eurovision-voting/internal/storage"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service struct {
	storage *storage.Storage
	jwt     *jwt.Service
	s3      *media.S3

	mu        sync.RWMutex
	userConns map[uuid.UUID]*websocket.Conn
}

func New(storage *storage.Storage, jwtSvc *jwt.Service, s3 *media.S3) *Service {
	return &Service{
		storage: storage,
		jwt:     jwtSvc,
		s3:      s3,
		userConns: make(map[uuid.UUID]*websocket.Conn),
	}
}
