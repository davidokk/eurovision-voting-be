package service

import (
	"eurovision-voting/internal/jwt"
	"eurovision-voting/internal/storage"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service struct {
	storage *storage.Storage
	jwt     *jwt.Service

	mu        sync.RWMutex
	userConns map[uuid.UUID]*websocket.Conn
}

func New(storage *storage.Storage, jwtSvc *jwt.Service) *Service {
	return &Service{
		storage: storage,
		jwt:     jwtSvc,

		userConns: make(map[uuid.UUID]*websocket.Conn),
	}
}
