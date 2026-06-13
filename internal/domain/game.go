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
	Custom        bool   `json:"custom,omitempty"`
}

type GamePlayer struct {
	UserID    string  `json:"user_id"`
	Username  string  `json:"username"`
	Score     int     `json:"score"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type GameContestScore struct {
	Username string  `json:"username"`
	Score    int     `json:"score"`
	Comment  *string `json:"comment,omitempty"`
}

type GameRoundView struct {
	Index         int     `json:"index"`
	PerformanceID string  `json:"performance_id"`
	YoutubeLink   string  `json:"youtube_link"`
	Mode          string  `json:"mode"` // audio | video | video_full | silent
	VideoStartSec int     `json:"video_start_sec,omitempty"`
	Artist        *string `json:"artist,omitempty"`
	Song          *string `json:"song,omitempty"`
	CountryName   *string `json:"country_name,omitempty"`
	FlagEmoji     *string `json:"flag_emoji,omitempty"`
	Year          *int    `json:"year,omitempty"`
	ContestType   *string `json:"contest_type,omitempty"`
	RoundEndsAt   *time.Time `json:"round_ends_at,omitempty"`
	ContestScores []GameContestScore `json:"contest_scores,omitempty"`
}

type GamePlaylistEntryInput struct {
	PerformanceID string `json:"performance_id"`
	Artist        string `json:"artist,omitempty"`
	Song          string `json:"song,omitempty"`
	YoutubeLink   string `json:"youtube_link,omitempty"`
}

type GameRoomView struct {
	Code             string             `json:"code"`
	HostUserID       string             `json:"host_user_id"`
	HostUsername     string             `json:"host_username"`
	State            string             `json:"state"`
	Paused           bool               `json:"paused"`
	PointsPerCorrect int                `json:"points_per_correct"`
	RoundDurationSec int                `json:"round_duration_sec"`
	Players          []GamePlayer       `json:"players"`
	PlaylistIDs      []string           `json:"playlist_ids,omitempty"`
	PlaylistPreview  []GameCatalogItem  `json:"playlist_preview,omitempty"`
	PlaylistMode     string             `json:"playlist_mode,omitempty"` // manual | auto
	AutoCount        int                `json:"auto_count,omitempty"`
	CurrentRound     int                `json:"current_round"`
	TotalRounds      int                `json:"total_rounds"`
	BuzzedUserID     *string            `json:"buzzed_user_id,omitempty"`
	BuzzedUsername   *string            `json:"buzzed_username,omitempty"`
	BuzzedAnswer     *string            `json:"buzzed_answer,omitempty"`
	PlayMode         string             `json:"play_mode,omitempty"` // offline | online
	Round            *GameRoundView     `json:"round,omitempty"`
	LastJudgement    *GameJudgement     `json:"last_judgement,omitempty"`
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
