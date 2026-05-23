package service

import (
	"eurovision-voting/internal/jwt"
	"eurovision-voting/internal/mail"
	"eurovision-voting/internal/media"
	"eurovision-voting/internal/storage"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const legacyPendingEmailSuffix = "@legacy.pending"

type Service struct {
	storage             *storage.Storage
	jwt                 *jwt.Service
	s3                  *media.S3
	mail                *mail.Client
	telegramBotUsername string
	signupAllowed       bool

	mu        sync.RWMutex
	userConns map[uuid.UUID]*websocket.Conn
}

func New(storage *storage.Storage, jwtSvc *jwt.Service, s3 *media.S3, mailClient *mail.Client, telegramBotUsername string, signupAllowed bool) *Service {
	return &Service{
		storage:             storage,
		jwt:                 jwtSvc,
		s3:                  s3,
		mail:                mailClient,
		telegramBotUsername: strings.TrimPrefix(strings.TrimSpace(telegramBotUsername), "@"),
		signupAllowed:       signupAllowed,
		userConns:           make(map[uuid.UUID]*websocket.Conn),
	}
}

func (s *Service) SetTelegramBotUsername(username string) {
	u := strings.TrimPrefix(strings.TrimSpace(username), "@")
	if u != "" {
		s.telegramBotUsername = u
	}
}
