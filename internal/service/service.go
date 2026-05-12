package service

import (
	"eurovision-voting/internal/jwt"
	"eurovision-voting/internal/storage"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service struct {
	storage *storage.Storage
	jwt     *jwt.Service

	userConns map[uuid.UUID]*websocket.Conn
}

func New(storage *storage.Storage, jwt *jwt.Service) *Service {
	return &Service{
		storage: storage,
		jwt:     jwt,

		userConns: make(map[uuid.UUID]*websocket.Conn),
	}
}
