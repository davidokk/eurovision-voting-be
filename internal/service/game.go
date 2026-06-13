package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	gameRevealDuration = 12 * time.Second
	gameCodeLen        = 6
)

var (
	ErrGameNotFound      = errors.New("game room not found")
	ErrGameForbidden     = errors.New("forbidden")
	ErrGameInvalidState  = errors.New("invalid game state")
	ErrGameAlreadyBuzzed = errors.New("someone already buzzed")
	ErrGameNotInRoom     = errors.New("not in room")
)

type gameRoomInternal struct {
	Code             string
	HostUserID       uuid.UUID
	HostUsername     string
	State            string
	Paused           bool
	PointsPerCorrect int
	Players          map[uuid.UUID]*domain.GamePlayer
	PlaylistIDs      []string
	Playlist         []domain.GameCatalogItem
	CurrentRound     int
	BuzzedUserID     *uuid.UUID
	BuzzedUsername   string
	RoundMode        string
	LastJudgement    *domain.GameJudgement
	revealTimer      *time.Timer
}

func randomGameCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, gameCodeLen)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}

func (s *Service) GetGameCatalog(ctx context.Context) ([]domain.GameCatalogItem, error) {
	return s.storage.ListGameCatalog(ctx)
}

func (s *Service) CreateGameRoom(ctx context.Context, hostID uuid.UUID) (*domain.GameRoomView, error) {
	user, err := s.storage.GetUser(ctx, hostID)
	if err != nil {
		return nil, err
	}

	var code string
	for i := 0; i < 20; i++ {
		code, err = randomGameCode()
		if err != nil {
			return nil, err
		}
		s.gameMu.RLock()
		_, exists := s.gameRooms[code]
		s.gameMu.RUnlock()
		if !exists {
			break
		}
	}

	room := &gameRoomInternal{
		Code:             code,
		HostUserID:       hostID,
		HostUsername:     user.Username,
		State:            "lobby",
		PointsPerCorrect: 10,
		Players:          make(map[uuid.UUID]*domain.GamePlayer),
		RoundMode:        "audio",
	}
	room.Players[hostID] = &domain.GamePlayer{
		UserID:   hostID.String(),
		Username: user.Username,
		Score:    0,
	}

	s.gameMu.Lock()
	s.gameRooms[code] = room
	s.gameMu.Unlock()

	return s.roomView(room), nil
}

func (s *Service) JoinGameRoom(ctx context.Context, code string, userID uuid.UUID) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	user, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.State == "finished" {
		return nil, ErrGameInvalidState
	}

	if _, exists := room.Players[userID]; !exists {
		room.Players[userID] = &domain.GamePlayer{
			UserID:   userID.String(),
			Username: user.Username,
			Score:    0,
		}
	}

	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) GetGameRoom(code string) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	s.gameMu.RLock()
	defer s.gameMu.RUnlock()
	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	return s.roomViewLocked(room), nil
}

func (s *Service) SetGamePlaylist(ctx context.Context, code string, hostID uuid.UUID, ids []string) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	items, err := s.storage.GetGameCatalogItems(ctx, ids)
	if err != nil {
		return nil, err
	}

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.HostUserID != hostID {
		return nil, ErrGameForbidden
	}
	if room.State != "lobby" {
		return nil, ErrGameInvalidState
	}

	room.PlaylistIDs = ids
	room.Playlist = items

	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) SetGamePoints(code string, hostID uuid.UUID, points int) (*domain.GameRoomView, error) {
	if points < 1 || points > 100 {
		points = 10
	}
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.HostUserID != hostID {
		return nil, ErrGameForbidden
	}
	room.PointsPerCorrect = points

	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) StartGame(code string, hostID uuid.UUID) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.HostUserID != hostID {
		return nil, ErrGameForbidden
	}
	if room.State != "lobby" || len(room.Playlist) == 0 {
		return nil, ErrGameInvalidState
	}

	room.CurrentRound = 0
	room.Paused = false
	room.LastJudgement = nil
	s.startRoundLocked(room)

	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) startRoundLocked(room *gameRoomInternal) {
	if room.revealTimer != nil {
		room.revealTimer.Stop()
		room.revealTimer = nil
	}
	room.BuzzedUserID = nil
	room.BuzzedUsername = ""
	room.LastJudgement = nil
	room.RoundMode = "audio"
	room.State = "round_playing"
}

func (s *Service) GameBuzz(code string, userID uuid.UUID) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if _, in := room.Players[userID]; !in {
		return nil, ErrGameNotInRoom
	}
	if room.Paused || room.State != "round_playing" {
		return nil, ErrGameInvalidState
	}
	if room.BuzzedUserID != nil {
		return nil, ErrGameAlreadyBuzzed
	}

	player := room.Players[userID]
	room.BuzzedUserID = &userID
	room.BuzzedUsername = player.Username
	room.State = "round_buzzed"

	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) GameJudge(code string, hostID uuid.UUID, correct bool) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.HostUserID != hostID {
		return nil, ErrGameForbidden
	}
	if room.State != "round_buzzed" || room.BuzzedUserID == nil {
		return nil, ErrGameInvalidState
	}

	buzzedID := *room.BuzzedUserID
	player := room.Players[buzzedID]
	delta := 0
	if correct {
		delta = room.PointsPerCorrect
		player.Score += delta
	}

	room.LastJudgement = &domain.GameJudgement{
		Correct:  correct,
		Username: player.Username,
		Points:   player.Score,
		Delta:    delta,
	}

	room.State = "round_reveal"
	room.RoundMode = "video"
	revealUntil := time.Now().Add(gameRevealDuration)
	roomCode := code

	if room.revealTimer != nil {
		room.revealTimer.Stop()
	}
	room.revealTimer = time.AfterFunc(gameRevealDuration, func() {
		s.onRevealEnd(roomCode)
	})

	view := s.roomViewLocked(room)
	view.Round.RevealUntil = &revealUntil
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) onRevealEnd(code string) {
	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok || room.State != "round_reveal" {
		return
	}
	room.revealTimer = nil

	if room.CurrentRound+1 >= len(room.Playlist) {
		room.State = "finished"
		view := s.roomViewLocked(room)
		s.broadcastGameRoomLocked(code, view)
		return
	}

	room.CurrentRound++
	s.startRoundLocked(room)
	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
}

func (s *Service) GameSkipToNext(code string, hostID uuid.UUID) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.HostUserID != hostID {
		return nil, ErrGameForbidden
	}
	if room.State != "round_reveal" {
		return nil, ErrGameInvalidState
	}

	if room.revealTimer != nil {
		room.revealTimer.Stop()
		room.revealTimer = nil
	}

	if room.CurrentRound+1 >= len(room.Playlist) {
		room.State = "finished"
	} else {
		room.CurrentRound++
		s.startRoundLocked(room)
	}

	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) GamePause(code string, hostID uuid.UUID, paused bool) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if room.HostUserID != hostID {
		return nil, ErrGameForbidden
	}
	if room.State == "lobby" || room.State == "finished" {
		return nil, ErrGameInvalidState
	}

	room.Paused = paused
	view := s.roomViewLocked(room)
	s.broadcastGameRoomLocked(code, view)
	return view, nil
}

func (s *Service) roomView(room *gameRoomInternal) *domain.GameRoomView {
	s.gameMu.RLock()
	defer s.gameMu.RUnlock()
	return s.roomViewLocked(room)
}

func (s *Service) roomViewLocked(room *gameRoomInternal) *domain.GameRoomView {
	players := make([]domain.GamePlayer, 0, len(room.Players))
	for _, p := range room.Players {
		players = append(players, *p)
	}

	view := &domain.GameRoomView{
		Code:             room.Code,
		HostUserID:       room.HostUserID.String(),
		HostUsername:     room.HostUsername,
		State:            room.State,
		Paused:           room.Paused,
		PointsPerCorrect: room.PointsPerCorrect,
		Players:          players,
		PlaylistIDs:      room.PlaylistIDs,
		CurrentRound:     room.CurrentRound,
		TotalRounds:      len(room.Playlist),
		LastJudgement:    room.LastJudgement,
	}

	if room.BuzzedUserID != nil {
		id := room.BuzzedUserID.String()
		view.BuzzedUserID = &id
		view.BuzzedUsername = &room.BuzzedUsername
	}

	if room.State != "lobby" && len(room.Playlist) > 0 && room.CurrentRound < len(room.Playlist) {
		item := room.Playlist[room.CurrentRound]
		round := &domain.GameRoundView{
			Index:         room.CurrentRound,
			PerformanceID: item.PerformanceID,
			YoutubeLink:   item.YoutubeLink,
			Mode:          room.RoundMode,
		}
		if room.State == "round_reveal" || room.State == "finished" {
			round.Artist = &item.Artist
			round.Song = &item.Song
			round.CountryName = &item.CountryName
			round.FlagEmoji = &item.FlagEmoji
			round.Year = &item.Year
			round.ContestType = &item.ContestType
		}
		view.Round = round
	}

	return view
}

func (s *Service) broadcastGameRoomLocked(code string, room *domain.GameRoomView) {
	conns := s.gameRoomConns[code]
	if len(conns) == 0 {
		return
	}

	event := domain.GameEvent{Type: "game.state", Room: room}
	var stale []uuid.UUID
	for userID, conn := range conns {
		_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		if err := conn.WriteJSON(event); err != nil {
			stale = append(stale, userID)
			_ = conn.Close()
		}
	}
	for _, id := range stale {
		delete(conns, id)
	}
}

func (s *Service) ServeGameConn(code string, userID uuid.UUID, username string, conn *websocket.Conn) {
	code = strings.ToUpper(strings.TrimSpace(code))

	s.gameMu.Lock()
	if s.gameRoomConns[code] == nil {
		s.gameRoomConns[code] = make(map[uuid.UUID]*websocket.Conn)
	}
	if old, ok := s.gameRoomConns[code][userID]; ok && old != nil {
		_ = old.Close()
	}
	s.gameRoomConns[code][userID] = conn

	room, roomExists := s.gameRooms[code]
	if roomExists {
		if _, in := room.Players[userID]; !in {
			room.Players[userID] = &domain.GamePlayer{
				UserID:   userID.String(),
				Username: username,
				Score:    0,
			}
		}
		view := s.roomViewLocked(room)
		s.gameMu.Unlock()

		_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		_ = conn.WriteJSON(domain.GameEvent{Type: "game.state", Room: view})
	} else {
		s.gameMu.Unlock()
		_ = conn.WriteJSON(map[string]string{"type": "game.error", "message": "room not found"})
		_ = conn.Close()
		return
	}

	conn.SetReadLimit(wsReadLimit)
	_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	go s.gameWritePump(code, userID, conn)
	go s.gameReadPump(code, userID, conn)
}

func (s *Service) gameWritePump(code string, userID uuid.UUID, conn *websocket.Conn) {
	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()

	for range ticker.C {
		s.gameMu.RLock()
		cur, ok := s.gameRoomConns[code][userID]
		s.gameMu.RUnlock()
		if !ok || cur != conn {
			return
		}
		_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			_ = conn.Close()
			return
		}
	}
}

func (s *Service) gameReadPump(code string, userID uuid.UUID, conn *websocket.Conn) {
	defer func() {
		_ = conn.Close()
		s.gameMu.Lock()
		if m, ok := s.gameRoomConns[code]; ok {
			if cur, ok := m[userID]; ok && cur == conn {
				delete(m, userID)
			}
		}
		s.gameMu.Unlock()
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))
		if messageType != websocket.TextMessage {
			continue
		}

		var m map[string]any
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		t, _ := m["type"].(string)
		switch t {
		case "ping":
			_ = conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			_ = conn.WriteJSON(map[string]string{"type": "pong"})
		case "game.buzz":
			if _, err := s.GameBuzz(code, userID); err != nil {
				log.Debug().Err(err).Str("code", code).Msg("game buzz")
			}
		default:
			// host actions via REST only for simplicity
		}
	}
}

func (s *Service) HandleGameHostAction(ctx context.Context, code string, hostID uuid.UUID, action string, payload map[string]any) (*domain.GameRoomView, error) {
	switch action {
	case "start":
		return s.StartGame(code, hostID)
	case "judge":
		correct, _ := payload["correct"].(bool)
		return s.GameJudge(code, hostID, correct)
	case "next":
		return s.GameSkipToNext(code, hostID)
	case "pause":
		return s.GamePause(code, hostID, true)
	case "resume":
		return s.GamePause(code, hostID, false)
	case "set_points":
		points := 10
		if v, ok := payload["points"].(float64); ok {
			points = int(v)
		}
		return s.SetGamePoints(code, hostID, points)
	case "set_playlist":
		var ids []string
		if raw, ok := payload["performance_ids"].([]any); ok {
			for _, x := range raw {
				if str, ok := x.(string); ok {
					ids = append(ids, str)
				}
			}
		}
		return s.SetGamePlaylist(ctx, code, hostID, ids)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
