package service

import (
	"context"
	"encoding/json"
	"eurovision-voting/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	wsReadLimit   = 512 * 1024
	wsPongWait    = 90 * time.Second
	wsPingPeriod  = 30 * time.Second
	wsWriteWait   = 10 * time.Second
)

// ServeConn регистрирует соединение, read/write pumps, ping/pong.
func (s *Service) ServeConn(userID uuid.UUID, conn *websocket.Conn) {
	s.mu.Lock()
	if old, ok := s.userConns[userID]; ok && old != nil {
		_ = old.Close()
	}
	s.userConns[userID] = conn
	s.mu.Unlock()

	conn.SetReadLimit(wsReadLimit)
	_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	go s.writePump(userID, conn)
	go s.readPump(userID, conn)
}

func (s *Service) writePump(userID uuid.UUID, conn *websocket.Conn) {
	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.RLock()
		cur, ok := s.userConns[userID]
		s.mu.RUnlock()
		if !ok || cur != conn {
			return
		}

		_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Debug().Err(err).Str("user", userID.String()).Msg("ws server ping failed")
			_ = conn.Close()
			return
		}
	}
}

func (s *Service) readPump(userID uuid.UUID, conn *websocket.Conn) {
	defer func() {
		_ = conn.Close()
		s.mu.Lock()
		if cur, ok := s.userConns[userID]; ok && cur == conn {
			delete(s.userConns, userID)
		}
		s.mu.Unlock()
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Debug().Err(err).Str("user", userID.String()).Msg("ws read end")
			}
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))

		if messageType == websocket.TextMessage {
			var m map[string]any
			if json.Unmarshal(data, &m) != nil {
				continue
			}
			if t, _ := m["type"].(string); t == "ping" {
				_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
				if err := conn.WriteJSON(map[string]string{"type": "pong"}); err != nil {
					return
				}
			}
		}
	}
}

func (s *Service) broadcastMessage(msg *domain.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var stale []uuid.UUID
	for userID, conn := range s.userConns {
		if err := conn.WriteJSON(msg); err != nil {
			stale = append(stale, userID)
			if websocket.IsCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure) {
				log.Debug().Err(err).Str("user", userID.String()).Msg("conn closed on broadcast")
			} else {
				log.Error().Err(err).Str("user", userID.String()).Msg("broadcast message")
			}
			_ = conn.Close()
		}
	}
	for _, id := range stale {
		delete(s.userConns, id)
	}
}

func (s *Service) SendMessage(ctx context.Context, msg *domain.Message) error {
	user, err := s.storage.GetUser(ctx, msg.UserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	msg.Username = user.Username
	msg.AvatarURL = user.AvatarURL
	msg.CreatedAt = time.Now()
	if msg.ContentType == "" {
		msg.ContentType = "text"
	}
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	if err := s.storage.InsertMessage(ctx, msg); err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	if err := s.storage.FillReplyPreview(ctx, msg); err != nil {
		log.Warn().Err(err).Str("message", msg.ID.String()).Msg("reply preview")
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
