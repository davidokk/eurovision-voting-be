package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Username       string
	HashedPassword string
	Role           *string
}

type Country struct {
	ID        string `json:"id"`
	NameRU    string `json:"name_ru"`
	FlagEmoji string `json:"flag_emoji"`
}

type Contest struct {
	ID     string    `json:"id"`
	Type   string    `json:"type"`
	Year   int       `json:"year"`
	Starts time.Time `json:"starts"`
	Ends   time.Time `json:"ends"`
}

type Performance struct {
	ID          string `json:"id"`
	ContestID   string `json:"contest_id"`
	CountryID   string `json:"country_id"`
	Number      int    `json:"number"`
	YoutubeLink string `json:"youtube_link"`
	Artist      string `json:"artist"`
	Song        string `json:"song"`
}

type Score struct {
	UserID        string  `json:"user_id"`
	PerformanceID string  `json:"performance_id"`
	Score         int     `json:"score"`
	Comment       *string `json:"comment,omitempty"`
}

type ScoreView struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Score    int     `json:"score"`
	Comment  *string `json:"comment,omitempty"`
}

type PerformanceWithScores struct {
	PerformanceID string      `json:"performance_id"`
	Country       Country     `json:"country"`
	Artist        string      `json:"artist"`
	Song          string      `json:"song"`
	Number        int         `json:"number"`
	YoutubeLink   string      `json:"youtube_link"`
	Scores        []ScoreView `json:"scores"`
	TotalScore    float64     `json:"total_score"`
	Qualified     bool        `json:"qualified"`
}

type ContestParticipantView struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type ContestView struct {
	Contest Contest `json:"contest"`
	// Participants []ContestParticipantView `json:"participants"`
	Performances []PerformanceWithScores `json:"performances"`
}
