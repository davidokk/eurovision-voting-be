package request

import (
	"encoding/json"
	"errors"
	"eurovision-voting/internal/domain"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type TelegramConfirmRequest struct {
	LinkToken string `json:"link_token"`
	Code      string `json:"code"`
}

func (r *TelegramConfirmRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *TelegramConfirmRequest) Validate() error {
	if strings.TrimSpace(r.LinkToken) == "" {
		return fmt.Errorf("link_token is empty")
	}
	if strings.TrimSpace(r.Code) == "" {
		return fmt.Errorf("code is empty")
	}
	return nil
}

type ChangeUsernameRequest struct {
	Username string `json:"username"`
}

func (r *ChangeUsernameRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *ChangeUsernameRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return fmt.Errorf("username is empty")
	}
	return nil
}

type RatePerformanceRequest struct {
	PerformanceID uuid.UUID
	Score         int    `json:"score"`
	Comment       string `json:"comment"`
	Gif           string `json:"gif_url"`
}

func (r *RatePerformanceRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return err
	}
	id := chi.URLParam(req, "id")
	r.PerformanceID, _ = uuid.Parse(id)
	return nil
}

func (r *RatePerformanceRequest) Validate() error {
	if r.PerformanceID == uuid.Nil {
		return errors.New("performance_id is missing")
	}
	if r.Score < 1 || r.Score > 10 {
		return errors.New("score must be between 1 and 10")
	}
	return nil
}

type UpdatePerformanceRequest struct {
	PerformanceID uuid.UUID
	Qualified     bool   `json:"qualified"`
	YoutubeLink   string `json:"youtube_link"`
	Place         *int   `json:"place"`
}

func (r *UpdatePerformanceRequest) Bind(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return err
	}
	id := chi.URLParam(req, "id")
	r.PerformanceID, _ = uuid.Parse(id)
	return nil
}

func (r *UpdatePerformanceRequest) Validate() error {
	if r.PerformanceID == uuid.Nil {
		return errors.New("performance_id is missing")
	}
	return nil
}

type GetScoresFilteredRequest struct {
	UserID      *uuid.UUID
	CountryID   *uuid.UUID
	ContestYear *int
	Sort        domain.SortType
}

func (r *GetScoresFilteredRequest) Bind(req *http.Request) error {
	q := req.URL.Query()

	if v := q.Get("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid user_id")
		}
		r.UserID = &id
	}

	if v := q.Get("country_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid country_id")
		}
		r.CountryID = &id
	}

	if v := q.Get("year"); v != "" {
		year, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid year")
		}
		r.ContestYear = &year
	}

	switch q.Get("sort") {
	case "score":
		r.Sort = domain.SortByScore
	case "time", "":
		r.Sort = domain.SortByTime
	default:
		return fmt.Errorf("invalid sort value")
	}

	return nil
}

func (r *GetScoresFilteredRequest) Validate() error {
	switch r.Sort {
	case domain.SortByTime,
		domain.SortByScore:
	default:
		return fmt.Errorf("invalid sort type")
	}

	return nil
}

type SendMessageRequest struct {
	ContestID       uuid.UUID  `json:"contest_id"`
	Message         string     `json:"message"`
	ContentType     string     `json:"content_type"`
	MediaURL        *string    `json:"media_url"`
	MediaDurationMs *int       `json:"media_duration_ms"`
	ReplyToID       *uuid.UUID `json:"reply_to_id"`
}

func (r *SendMessageRequest) Bind(req *http.Request) error {
	if req.Header.Get("Content-Type") == "application/json" || req.ContentLength > 0 {
		if err := json.NewDecoder(req.Body).Decode(r); err != nil {
			return err
		}
		return nil
	}
	ci := req.URL.Query().Get("contest_id")
	var err error
	r.ContestID, err = uuid.Parse(ci)
	if err != nil {
		return fmt.Errorf("contest_id: %w", err)
	}
	r.Message = req.URL.Query().Get("message")
	r.ContentType = "text"
	return nil
}

func (r *SendMessageRequest) Validate() error {
	if r.ContestID == uuid.Nil {
		return fmt.Errorf("contest_id required")
	}
	ct := r.ContentType
	if ct == "" {
		ct = "text"
		r.ContentType = ct
	}
	switch ct {
	case "text":
		if strings.TrimSpace(r.Message) == "" {
			return fmt.Errorf("message empty")
		}
	case "voice", "video_note", "image":
		if r.MediaURL == nil || *r.MediaURL == "" {
			return fmt.Errorf("media_url required")
		}
	default:
		return fmt.Errorf("invalid content_type")
	}
	return nil
}

type AdminCreateContestRequest struct {
	Year  int    `json:"year"`
	Type  string `json:"type"`
	Starts string `json:"starts"`
	Ends   string `json:"ends"`
}

func (r *AdminCreateContestRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *AdminCreateContestRequest) Validate() error {
	if r.Year < 2000 || r.Year > 2100 {
		return fmt.Errorf("invalid year")
	}
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("type required")
	}
	if strings.TrimSpace(r.Starts) == "" || strings.TrimSpace(r.Ends) == "" {
		return fmt.Errorf("starts and ends required")
	}
	return nil
}

type AdminCreatePerformanceRequest struct {
	CountryID    string `json:"country_id"`
	Number       int    `json:"number"`
	Artist       string `json:"artist"`
	Song         string `json:"song"`
	YoutubeLink  string `json:"youtube_link"`
}

func (r *AdminCreatePerformanceRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *AdminCreatePerformanceRequest) Validate() error {
	if _, err := uuid.Parse(r.CountryID); err != nil {
		return fmt.Errorf("invalid country_id")
	}
	if strings.TrimSpace(r.Artist) == "" || strings.TrimSpace(r.Song) == "" {
		return fmt.Errorf("artist and song required")
	}
	return nil
}

type AdminUpdatePlacesRequest struct {
	Ranked []string `json:"ranked"`
}

func (r *AdminUpdatePlacesRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *AdminUpdatePlacesRequest) Validate() error {
	return nil
}
