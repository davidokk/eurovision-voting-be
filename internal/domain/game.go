package domain

import "time"

type GameCatalogItem struct {
	PerformanceID string `json:"performance_id"`
	Artist        string `json:"artist"`
	Song          string `json:"song"`
	CountryName   string `json:"country_name"`
	FlagEmoji     string `json:"flag_emoji"`
	Year          int    `json:"year"`
	ContestType   string `json:"contest_type"`
	YoutubeLink   string `json:"youtube_link"`
}

type GamePlayer struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type GameRoundView struct {
	Index         int     `json:"index"`
	PerformanceID string  `json:"performance_id"`
	YoutubeLink   string  `json:"youtube_link"`
	Mode          string  `json:"mode"` // audio | video
	Artist        *string `json:"artist,omitempty"`
	Song          *string `json:"song,omitempty"`
	CountryName   *string `json:"country_name,omitempty"`
	FlagEmoji     *string `json:"flag_emoji,omitempty"`
	Year          *int    `json:"year,omitempty"`
	ContestType   *string `json:"contest_type,omitempty"`
	RevealUntil   *time.Time `json:"reveal_until,omitempty"`
}

type GameRoomView struct {
	Code             string           `json:"code"`
	HostUserID       string           `json:"host_user_id"`
	HostUsername     string           `json:"host_username"`
	State            string           `json:"state"`
	Paused           bool             `json:"paused"`
	PointsPerCorrect int              `json:"points_per_correct"`
	Players          []GamePlayer     `json:"players"`
	PlaylistIDs      []string         `json:"playlist_ids,omitempty"`
	CurrentRound     int              `json:"current_round"`
	TotalRounds      int              `json:"total_rounds"`
	BuzzedUserID     *string          `json:"buzzed_user_id,omitempty"`
	BuzzedUsername   *string          `json:"buzzed_username,omitempty"`
	Round            *GameRoundView   `json:"round,omitempty"`
	LastJudgement    *GameJudgement   `json:"last_judgement,omitempty"`
}

type GameJudgement struct {
	Correct  bool   `json:"correct"`
	Username string `json:"username"`
	Points   int    `json:"points"`
	Delta    int    `json:"delta"`
}

type GameEvent struct {
	Type string       `json:"type"`
	Room *GameRoomView `json:"room,omitempty"`
}
