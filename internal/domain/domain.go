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
	ID     uuid.UUID `json:"id"`
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
	GifURL   *string `json:"gif_url,omitempty"`
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
	Place         *int        `json:"place"`
}

type ContestParticipantView struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type ContestView struct {
	Contest      Contest                 `json:"contest"`
	Performances []PerformanceWithScores `json:"performances"`
}

type ScoreFiltered struct {
	Username    string
	CountryName string
	ContestYear int
	ContestType string
	Score       int
	Comment     *string
	YoutubeLink string
	GifURL      *string
	Song        string
	Artist      string
	Qualified   *bool
	Place       *int
}

type SortType int

const (
	SortByTime SortType = iota
	SortByScore
)

type Filters struct {
	UserID      *uuid.UUID
	CountryID   *uuid.UUID
	ContestYear *int
	Sort        SortType
}

type Message struct {
	UserID        uuid.UUID
	PerformanceID uuid.UUID
	ContestID     uuid.UUID `json:"contestId"`
	Username      string    `json:"username"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"createdAt"`
	Type          string    `json:"type"`
	Gif           *string   `json:"gif"`
	Country       *string   `json:"country"`
	CountryFlag   *string   `json:"country_flag"`
	Score         *int      `json:"score"`
	OldScore      *int      `json:"old_score"`
	Comment       *string   `json:"comment"`
}
