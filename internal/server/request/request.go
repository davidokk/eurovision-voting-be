package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	Qualified     bool `json:"qualified"`
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
