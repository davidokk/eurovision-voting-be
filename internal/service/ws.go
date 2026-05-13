package service

import (
	"context"
	"eurovision-voting/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func (s *Service) ServeConn(userID uuid.UUID, conn *websocket.Conn) {
	s.userConns[userID] = conn
}

func (s *Service) broadcastMessage(msg *domain.Message) {
	log.Info().Int("users", len(s.userConns)).Msg("broadcast")
	for userID, conn := range s.userConns {
		if err := conn.WriteJSON(msg); err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure) {
				delete(s.userConns, userID)
				log.Error().Err(err).Msg("conn closed")
			} else {
				log.Error().Err(err).Msg("broadcast message")
			}
		}
		log.Info().Msg("msg sent")
	}
}

func (s *Service) SendMessage(ctx context.Context, msg *domain.Message) error {
	user, err := s.storage.GetUser(ctx, msg.UserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	msg.Username = user.Username
	msg.CreatedAt = time.Now()
	if err := s.storage.InsertMessage(ctx, msg); err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	s.broadcastMessage(msg)
	return nil
}

func (s *Service) GetMessages(ctx context.Context, contestID uuid.UUID) ([]domain.Message, error) {
	res, err := s.storage.GetMessages(ctx, contestID)
	if err != nil {
		log.Error().Err(err).Send()
	}
	return res, err
}