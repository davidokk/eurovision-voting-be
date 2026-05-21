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

type SingUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *SingUpRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *SingUpRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if r.Password == "" {
		return fmt.Errorf("password is empty")
	}
	return nil
}

type SingInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *SingInRequest) Bind(req *http.Request) error {
	return json.NewDecoder(req.Body).Decode(r)
}

func (r *SingInRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if r.Password == "" {
		return fmt.Errorf("password is empty")
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
	Place         int    `json:"place"`
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
