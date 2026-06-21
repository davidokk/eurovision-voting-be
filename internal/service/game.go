package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"eurovision-voting/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	defaultGameRoundDurationSec = 10
	gameCodeLen                 = 6
	roundCountdownSec           = 3
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
	RoundDurationSec int
	Players          map[uuid.UUID]*domain.GamePlayer
	PlaylistIDs      []string
	Playlist         []domain.GameCatalogItem
	PlaylistMode     string
	AutoCount        int
	CurrentRound     int
	BuzzedUserID     *uuid.UUID
	BuzzedUsername   string
	BuzzedAnswer     string
	PlayMode         string
	RoundMode        string
	VideoStartSec    int
	RoundEndsAt      *time.Time
	ContestScores     []domain.GameContestScore
	ContestTotalScore float64
	ContestQualified  *bool
	ContestPlace      *int
	LastJudgement     *domain.GameJudgement
	roundTimer       *time.Timer
	clipTimer        *time.Timer
	roundPauseRemaining time.Duration
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

func randomVideoStartSec() int {
	n, err := rand.Int(rand.Reader, big.NewInt(90))
	if err != nil {
		return 45
	}
	return int(n.Int64()) + 15
}

func gamePlayerFromUser(user *domain.User, score int) *domain.GamePlayer {
	p := &domain.GamePlayer{
		UserID:   user.ID.String(),
		Username: user.Username,
		Score:    score,
	}
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		p.AvatarURL = user.AvatarURL
	}
	return p
}

func clampRoundDurationSec(sec int) int {
	if sec < 3 {
		return 3
	}
	if sec > 120 {
		return 120
	}
	return sec
}

func (s *Service) roundPlayDuration(room *gameRoomInternal) time.Duration {
	sec := room.RoundDurationSec
	if sec <= 0 {
		sec = defaultGameRoundDurationSec
	}
	return time.Duration(sec) * time.Second
}

func (s *Service) stopRoundTimer(room *gameRoomInternal) {
	if room.roundTimer != nil {
		room.roundTimer.Stop()
		room.roundTimer = nil
	}
}

func (s *Service) stopClipTimer(room *gameRoomInternal) {
	if room.clipTimer != nil {
		room.clipTimer.Stop()
		room.clipTimer = nil
	}
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
		RoundDurationSec: defaultGameRoundDurationSec,
		Players:          make(map[uuid.UUID]*domain.GamePlayer),
		RoundMode:        "audio",
		PlaylistMode:     "manual",
		AutoCount:        10,
		PlayMode:         "offline",
	}
	room.Players[hostID] = gamePlayerFromUser(user, 0)

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
		room.Players[userID] = gamePlayerFromUser(user, 0)
	} else if user.AvatarURL != nil && *user.AvatarURL != "" {
		room.Players[userID].AvatarURL = user.AvatarURL
	}

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, userID), nil
}

func (s *Service) GetGameRoom(code string) (*domain.GameRoomView, error) {
	return s.GetGameRoomForUser(code, uuid.Nil)
}

func (s *Service) GetGameRoomForUser(code string, userID uuid.UUID) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	s.gameMu.RLock()
	defer s.gameMu.RUnlock()
	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if userID != uuid.Nil {
		return s.roomViewForUserLocked(room, userID), nil
	}
	return s.roomViewLocked(room), nil
}

func isCustomTrackID(id string) bool {
	return strings.HasPrefix(id, "yt:")
}

var youtubeIDRe = regexp.MustCompile(`(?i)(?:youtu\.be/|youtube\.com.*v=)([^&?/]+)`)

func extractYouTubeVideoID(link string) string {
	m := youtubeIDRe.FindStringSubmatch(strings.TrimSpace(link))
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func customTrackItem(entry domain.GamePlaylistEntryInput) (domain.GameCatalogItem, error) {
	id := strings.TrimSpace(entry.PerformanceID)
	link := strings.TrimSpace(entry.YoutubeLink)
	videoID := ""
	if isCustomTrackID(id) {
		videoID = strings.TrimPrefix(id, "yt:")
	} else {
		videoID = extractYouTubeVideoID(link)
		id = "yt:" + videoID
	}
	if videoID == "" {
		return domain.GameCatalogItem{}, fmt.Errorf("invalid youtube track")
	}
	if link == "" {
		link = "https://www.youtube.com/watch?v=" + videoID
	}
	artist := strings.TrimSpace(entry.Artist)
	song := strings.TrimSpace(entry.Song)
	if artist == "" {
		artist = "YouTube"
	}
	if song == "" {
		song = videoID
	}
	return domain.GameCatalogItem{
		PerformanceID: id,
		Artist:        artist,
		Song:          song,
		CountryName:   "YouTube",
		FlagEmoji:     "🎵",
		YoutubeLink:   link,
		Custom:        true,
	}, nil
}

func (s *Service) buildPlaylistFromEntries(ctx context.Context, entries []domain.GamePlaylistEntryInput) ([]domain.GameCatalogItem, []string, error) {
	if len(entries) == 0 {
		return nil, nil, nil
	}

	var catalogIDs []string
	customEntries := make(map[string]domain.GamePlaylistEntryInput)
	orderedIDs := make([]string, 0, len(entries))

	for _, e := range entries {
		id := strings.TrimSpace(e.PerformanceID)
		link := strings.TrimSpace(e.YoutubeLink)

		if isCustomTrackID(id) {
			customEntries[id] = e
			orderedIDs = append(orderedIDs, id)
			continue
		}

		if id != "" {
			if _, err := uuid.Parse(id); err == nil {
				catalogIDs = append(catalogIDs, id)
				orderedIDs = append(orderedIDs, id)
				continue
			}
		}

		if link != "" {
			item, err := customTrackItem(e)
			if err != nil {
				continue
			}
			customEntries[item.PerformanceID] = e
			orderedIDs = append(orderedIDs, item.PerformanceID)
			continue
		}

		if id != "" {
			catalogIDs = append(catalogIDs, id)
			orderedIDs = append(orderedIDs, id)
		}
	}

	byID := make(map[string]domain.GameCatalogItem)
	if len(catalogIDs) > 0 {
		items, err := s.storage.GetGameCatalogItems(ctx, catalogIDs)
		if err != nil {
			return nil, nil, err
		}
		for _, item := range items {
			byID[item.PerformanceID] = item
		}
	}

	playlist := make([]domain.GameCatalogItem, 0, len(orderedIDs))
	ids := make([]string, 0, len(orderedIDs))
	seen := make(map[string]struct{})
	for _, id := range orderedIDs {
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		if isCustomTrackID(id) {
			entry, ok := customEntries[id]
			if !ok {
				continue
			}
			item, err := customTrackItem(entry)
			if err != nil {
				continue
			}
			playlist = append(playlist, item)
			ids = append(ids, item.PerformanceID)
			continue
		}
		if item, ok := byID[id]; ok {
			playlist = append(playlist, item)
			ids = append(ids, item.PerformanceID)
		}
	}
	return playlist, ids, nil
}

func shuffleGameCatalog(items []domain.GameCatalogItem, count int) []domain.GameCatalogItem {
	if count <= 0 || len(items) == 0 {
		return nil
	}
	if count > len(items) {
		count = len(items)
	}
	idx := make([]int, len(items))
	for i := range idx {
		idx[i] = i
	}
	for i := len(idx) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			j := i
			idx[i], idx[j] = idx[j], idx[i]
			continue
		}
		j := int(jBig.Int64())
		idx[i], idx[j] = idx[j], idx[i]
	}
	out := make([]domain.GameCatalogItem, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, items[idx[i]])
	}
	return out
}

func (s *Service) applyAutoPlaylistLocked(ctx context.Context, room *gameRoomInternal) error {
	catalog, err := s.storage.ListGameCatalog(ctx)
	if err != nil {
		return err
	}
	count := room.AutoCount
	if count < 1 {
		count = 10
	}
	picked := shuffleGameCatalog(catalog, count)
	if len(picked) == 0 {
		return ErrGameInvalidState
	}
	room.Playlist = picked
	room.PlaylistIDs = make([]string, len(picked))
	for i, item := range picked {
		room.PlaylistIDs[i] = item.PerformanceID
	}
	return nil
}

func parsePlaylistEntries(payload map[string]any) []domain.GamePlaylistEntryInput {
	if raw, ok := payload["entries"].([]any); ok && len(raw) > 0 {
		var entries []domain.GamePlaylistEntryInput
		for _, x := range raw {
			m, ok := x.(map[string]any)
			if !ok {
				continue
			}
			e := domain.GamePlaylistEntryInput{}
			if v, ok := m["performance_id"].(string); ok {
				e.PerformanceID = v
			}
			if v, ok := m["artist"].(string); ok {
				e.Artist = v
			}
			if v, ok := m["song"].(string); ok {
				e.Song = v
			}
			if v, ok := m["youtube_link"].(string); ok {
				e.YoutubeLink = v
			}
			if e.PerformanceID != "" || e.YoutubeLink != "" {
				entries = append(entries, e)
			}
		}
		return entries
	}
	if raw, ok := payload["performance_ids"].([]any); ok {
		var entries []domain.GamePlaylistEntryInput
		for _, x := range raw {
			if str, ok := x.(string); ok && str != "" {
				entries = append(entries, domain.GamePlaylistEntryInput{PerformanceID: str})
			}
		}
		return entries
	}
	return nil
}

func (s *Service) SetGamePlaylist(ctx context.Context, code string, hostID uuid.UUID, entries []domain.GamePlaylistEntryInput) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	playlist, ids, err := s.buildPlaylistFromEntries(ctx, entries)
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
	room.Playlist = playlist
	room.PlaylistMode = "manual"

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) SetGamePlaylistMode(code string, hostID uuid.UUID, mode string, count int) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if mode != "auto" {
		mode = "manual"
	}
	if count < 1 {
		count = 10
	}
	if count > 100 {
		count = 100
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

	room.PlaylistMode = mode
	room.AutoCount = count

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) SetGamePlaylistAuto(ctx context.Context, code string, hostID uuid.UUID, count int) (*domain.GameRoomView, error) {
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
	if room.State != "lobby" {
		return nil, ErrGameInvalidState
	}

	if count >= 1 && count <= 100 {
		room.AutoCount = count
	}
	room.PlaylistMode = "auto"
	if err := s.applyAutoPlaylistLocked(ctx, room); err != nil {
		return nil, err
	}

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) SetGameRoundDuration(code string, hostID uuid.UUID, seconds int) (*domain.GameRoomView, error) {
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
	if room.State != "lobby" {
		return nil, ErrGameInvalidState
	}
	room.RoundDurationSec = clampRoundDurationSec(seconds)

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) SetGamePlayMode(code string, hostID uuid.UUID, mode string) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if mode != "online" {
		mode = "offline"
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
	room.PlayMode = mode

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
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

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) StartGame(ctx context.Context, code string, hostID uuid.UUID) (*domain.GameRoomView, error) {
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
	if room.State != "lobby" {
		return nil, ErrGameInvalidState
	}
	if room.PlaylistMode == "auto" && len(room.Playlist) == 0 {
		if err := s.applyAutoPlaylistLocked(ctx, room); err != nil {
			return nil, err
		}
	}
	if len(room.Playlist) == 0 {
		return nil, ErrGameInvalidState
	}

	room.CurrentRound = 0
	room.Paused = false
	room.LastJudgement = nil
	s.startRoundLocked(code, room)

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) startRoundLocked(code string, room *gameRoomInternal) {
	s.stopRoundTimer(room)
	s.stopClipTimer(room)
	room.roundPauseRemaining = 0
	room.BuzzedUserID = nil
	room.BuzzedUsername = ""
	room.BuzzedAnswer = ""
	room.LastJudgement = nil
	room.ContestScores = nil
	room.ContestTotalScore = 0
	room.ContestQualified = nil
	room.ContestPlace = nil
	room.VideoStartSec = randomVideoStartSec()
	room.RoundMode = "silent"
	room.State = "round_countdown"
	ends := time.Now().Add(roundCountdownSec * time.Second)
	room.RoundEndsAt = &ends

	room.roundTimer = time.AfterFunc(roundCountdownSec*time.Second, func() {
		s.onRoundCountdownDone(code)
	})
}

func (s *Service) onRoundCountdownDone(code string) {
	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok || room.State != "round_countdown" {
		return
	}
	s.beginRoundPlayingLocked(code, room)
	s.broadcastGameRoomLocked(code, room)
}

func (s *Service) beginRoundPlayingLocked(code string, room *gameRoomInternal) {
	room.State = "round_playing"
	room.RoundMode = "audio"
	dur := s.roundPlayDuration(room)
	ends := time.Now().Add(dur)
	room.RoundEndsAt = &ends
	s.stopRoundTimer(room)
	room.roundTimer = time.AfterFunc(dur, func() {
		s.onRoundPlayTimeout(code)
	})
}

func (s *Service) onRoundPlayTimeout(code string) {
	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok || room.State != "round_playing" {
		return
	}
	s.stopRoundTimer(room)
	room.RoundEndsAt = nil
	s.enterRevealLocked(context.Background(), room)

	s.broadcastGameRoomLocked(code, room)
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

	s.stopRoundTimer(room)
	room.RoundEndsAt = nil

	player := room.Players[userID]
	room.BuzzedUserID = &userID
	room.BuzzedUsername = player.Username
	room.BuzzedAnswer = ""
	room.State = "round_buzzed"
	room.RoundMode = "silent"

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, userID), nil
}

const maxGameAnswerLen = 200

func (s *Service) GameSubmitAnswer(code string, userID uuid.UUID, answer string) (*domain.GameRoomView, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return nil, fmt.Errorf("answer required")
	}
	if len([]rune(answer)) > maxGameAnswerLen {
		return nil, fmt.Errorf("answer too long")
	}

	s.gameMu.Lock()
	defer s.gameMu.Unlock()

	room, ok := s.gameRooms[code]
	if !ok {
		return nil, ErrGameNotFound
	}
	if _, in := room.Players[userID]; !in {
		return nil, ErrGameNotInRoom
	}
	if room.PlayMode != "online" {
		return nil, ErrGameInvalidState
	}
	if room.State != "round_buzzed" || room.BuzzedUserID == nil {
		return nil, ErrGameInvalidState
	}
	if *room.BuzzedUserID != userID {
		return nil, ErrGameForbidden
	}
	if room.BuzzedAnswer != "" {
		return nil, fmt.Errorf("answer already submitted")
	}

	room.BuzzedAnswer = answer

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, userID), nil
}

func (s *Service) loadContestScores(ctx context.Context, room *gameRoomInternal) {
	if room.CurrentRound >= len(room.Playlist) {
		return
	}
	item := room.Playlist[room.CurrentRound]
	if item.Custom || isCustomTrackID(item.PerformanceID) {
		return
	}
	perfID, err := uuid.Parse(item.PerformanceID)
	if err != nil {
		return
	}
	scores, err := s.storage.GetScoresForPerformance(ctx, perfID)
	if err != nil {
		log.Debug().Err(err).Msg("game contest scores")
	} else {
		room.ContestScores = scores
	}
	stats, err := s.storage.GetPerformanceRevealStats(ctx, perfID)
	if err != nil {
		log.Debug().Err(err).Msg("game performance reveal stats")
		return
	}
	room.ContestTotalScore = stats.TotalScore
	if stats.Qualified.Valid {
		q := stats.Qualified.Bool
		room.ContestQualified = &q
	}
	if stats.Place.Valid && stats.Place.Int32 > 0 {
		p := int(stats.Place.Int32)
		room.ContestPlace = &p
	}
}

func (s *Service) enterRevealLocked(ctx context.Context, room *gameRoomInternal) {
	room.State = "round_reveal"
	room.RoundMode = "video"
	s.loadContestScores(ctx, room)
}

func gameJudgeDelta(outcome string, pointsPerCorrect int) (delta int, correct bool) {
	switch outcome {
	case "full":
		return pointsPerCorrect, true
	case "half":
		half := pointsPerCorrect / 2
		if half <= 0 {
			half = 1
		}
		return half, true
	default:
		return -pointsPerCorrect, false
	}
}

func parseJudgeOutcome(payload map[string]any) string {
	if o, ok := payload["outcome"].(string); ok {
		switch strings.TrimSpace(strings.ToLower(o)) {
		case "full", "half", "wrong":
			return strings.TrimSpace(strings.ToLower(o))
		}
	}
	if correct, ok := payload["correct"].(bool); ok {
		if correct {
			return "full"
		}
		return "wrong"
	}
	return "wrong"
}

func (s *Service) GameJudge(ctx context.Context, code string, hostID uuid.UUID, outcome string) (*domain.GameRoomView, error) {
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
	delta, correct := gameJudgeDelta(outcome, room.PointsPerCorrect)
	player.Score += delta

	room.LastJudgement = &domain.GameJudgement{
		Correct:  correct,
		Outcome:  outcome,
		Username: player.Username,
		Points:   player.Score,
		Delta:    delta,
	}

	s.enterRevealLocked(ctx, room)

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) GameRevealAnswer(ctx context.Context, code string, hostID uuid.UUID) (*domain.GameRoomView, error) {
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
	if room.State != "round_waiting_reveal" {
		return nil, ErrGameInvalidState
	}

	room.LastJudgement = nil
	s.enterRevealLocked(ctx, room)

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) GameStartFullClip(code string, hostID uuid.UUID) (*domain.GameRoomView, error) {
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

	s.stopClipTimer(room)
	room.State = "round_clip"
	room.RoundMode = "video_full"

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
}

func (s *Service) GameAdvanceRound(code string, hostID uuid.UUID) (*domain.GameRoomView, error) {
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
	if room.State != "round_reveal" && room.State != "round_clip" {
		return nil, ErrGameInvalidState
	}

	s.stopClipTimer(room)

	if room.CurrentRound+1 >= len(room.Playlist) {
		room.State = "finished"
		room.RoundMode = "silent"
	} else {
		room.CurrentRound++
		s.startRoundLocked(code, room)
	}

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
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

	if paused {
		room.Paused = true
		if room.State == "round_playing" && room.RoundEndsAt != nil {
			remaining := time.Until(*room.RoundEndsAt)
			if remaining < 0 {
				remaining = 0
			}
			s.stopRoundTimer(room)
			room.roundPauseRemaining = remaining
			snap := time.Now().Add(remaining)
			room.RoundEndsAt = &snap
		}
	} else {
		room.Paused = false
		if room.State == "round_playing" && room.roundPauseRemaining > 0 {
			remaining := room.roundPauseRemaining
			room.roundPauseRemaining = 0
			ends := time.Now().Add(remaining)
			room.RoundEndsAt = &ends
			s.stopRoundTimer(room)
			room.roundTimer = time.AfterFunc(remaining, func() {
				s.onRoundPlayTimeout(code)
			})
		}
	}

	s.broadcastGameRoomLocked(code, room)
	return s.roomViewForUserLocked(room, hostID), nil
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
		RoundDurationSec: room.RoundDurationSec,
		Players:          players,
		PlaylistIDs:      room.PlaylistIDs,
		PlaylistMode:     room.PlaylistMode,
		AutoCount:        room.AutoCount,
		CurrentRound:     room.CurrentRound,
		TotalRounds:      len(room.Playlist),
		LastJudgement:    room.LastJudgement,
	}

	if room.RoundDurationSec <= 0 {
		view.RoundDurationSec = defaultGameRoundDurationSec
	}

	if room.State == "lobby" && len(room.Playlist) > 0 {
		preview := make([]domain.GameCatalogItem, len(room.Playlist))
		copy(preview, room.Playlist)
		view.PlaylistPreview = preview
	}

	if len(room.Playlist) > 0 {
		sources := make([]string, 0, len(room.Playlist))
		for _, item := range room.Playlist {
			if item.YoutubeLink != "" {
				sources = append(sources, item.YoutubeLink)
			}
		}
		if len(sources) > 0 {
			view.PlaylistSources = sources
		}
	}

	if room.BuzzedUserID != nil {
		id := room.BuzzedUserID.String()
		view.BuzzedUserID = &id
		view.BuzzedUsername = &room.BuzzedUsername
	}
	if room.BuzzedAnswer != "" {
		ans := room.BuzzedAnswer
		view.BuzzedAnswer = &ans
	}
	if room.PlayMode != "" {
		view.PlayMode = room.PlayMode
	} else {
		view.PlayMode = "offline"
	}

	if round := s.buildRoundViewLocked(room, false); round != nil {
		view.Round = round
	}

	return view
}

func (s *Service) roomViewForUserLocked(room *gameRoomInternal, userID uuid.UUID) *domain.GameRoomView {
	forHost := userID != uuid.Nil && userID == room.HostUserID
	view := s.roomViewLocked(room)
	if forHost {
		if round := s.buildRoundViewLocked(room, true); round != nil {
			view.Round = round
		}
	}
	return view
}

func (s *Service) isActiveGameState(state string) bool {
	switch state {
	case "round_countdown", "round_playing", "round_buzzed", "round_waiting_reveal", "round_reveal", "round_clip":
		return true
	default:
		return false
	}
}

func (s *Service) buildRoundViewLocked(room *gameRoomInternal, forHost bool) *domain.GameRoundView {
	if !s.isActiveGameState(room.State) || len(room.Playlist) == 0 {
		return nil
	}

	idx := room.CurrentRound
	if idx < 0 {
		idx = 0
	}
	if idx >= len(room.Playlist) {
		idx = len(room.Playlist) - 1
	}

	item := room.Playlist[idx]
	round := &domain.GameRoundView{
		Index:         idx,
		PerformanceID: item.PerformanceID,
		YoutubeLink:   item.YoutubeLink,
		Mode:          room.RoundMode,
		VideoStartSec: room.VideoStartSec,
	}
	if room.RoundEndsAt != nil {
		round.RoundEndsAt = room.RoundEndsAt
	}
	revealed := room.State == "round_reveal" || room.State == "round_clip"
	hostHint := forHost && (room.State == "round_buzzed" || room.State == "round_waiting_reveal")
	if revealed || hostHint {
		round.Artist = &item.Artist
		round.Song = &item.Song
		round.CountryName = &item.CountryName
		round.FlagEmoji = &item.FlagEmoji
		round.Year = &item.Year
		round.ContestType = &item.ContestType
		if len(room.ContestScores) > 0 {
			round.ContestScores = room.ContestScores
		}
		if len(room.ContestScores) > 0 || room.ContestTotalScore > 0 {
			ts := room.ContestTotalScore
			round.TotalScore = &ts
		}
		if room.ContestQualified != nil {
			q := *room.ContestQualified
			round.Qualified = &q
		}
		if room.ContestPlace != nil {
			p := *room.ContestPlace
			round.Place = &p
		}
	}
	return round
}

func (s *Service) broadcastGameRoomLocked(code string, room *gameRoomInternal) {
	conns := s.gameRoomConns[code]
	if len(conns) == 0 {
		return
	}

	var stale []uuid.UUID
	for userID, conn := range conns {
		view := s.roomViewForUserLocked(room, userID)
		event := domain.GameEvent{Type: "game.state", Room: view}
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
		} else if room.Players[userID].Username == "" {
			room.Players[userID].Username = username
		}
		view := s.roomViewForUserLocked(room, userID)
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
		}
	}
}

func (s *Service) HandleGameHostAction(ctx context.Context, code string, hostID uuid.UUID, action string, payload map[string]any) (*domain.GameRoomView, error) {
	switch action {
	case "start":
		return s.StartGame(ctx, code, hostID)
	case "judge":
		outcome := parseJudgeOutcome(payload)
		return s.GameJudge(ctx, code, hostID, outcome)
	case "reveal":
		return s.GameRevealAnswer(ctx, code, hostID)
	case "clip":
		return s.GameStartFullClip(code, hostID)
	case "next_round":
		return s.GameAdvanceRound(code, hostID)
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
	case "set_round_duration":
		seconds := defaultGameRoundDurationSec
		if v, ok := payload["seconds"].(float64); ok {
			seconds = int(v)
		}
		return s.SetGameRoundDuration(code, hostID, seconds)
	case "set_playlist":
		entries := parsePlaylistEntries(payload)
		return s.SetGamePlaylist(ctx, code, hostID, entries)
	case "set_playlist_mode":
		mode, _ := payload["mode"].(string)
		count := 10
		if v, ok := payload["count"].(float64); ok {
			count = int(v)
		}
		return s.SetGamePlaylistMode(code, hostID, mode, count)
	case "set_playlist_auto":
		count := 0
		if v, ok := payload["count"].(float64); ok {
			count = int(v)
		}
		return s.SetGamePlaylistAuto(ctx, code, hostID, count)
	case "set_play_mode":
		mode, _ := payload["mode"].(string)
		return s.SetGamePlayMode(code, hostID, mode)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}
