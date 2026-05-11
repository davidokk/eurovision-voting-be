package server

import (
	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/server/request"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (s *Server) getContestViewHandler(w http.ResponseWriter, r *http.Request) {
	contestID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{
			Code: ValidatationCode,
			Err:  fmt.Sprintf("invalid or missing contest_id: %s", err.Error()),
		})
		return
	}

	res, err := s.service.GetContestView(r.Context(), contestID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	EncodeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) getContestsByYearHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.service.GetContestsByYears(r.Context())
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	EncodeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) ratePerformance(w http.ResponseWriter, r *http.Request) {
	var req request.RatePerformanceRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{
			Err:  err.Error(),
			Code: ValidatationCode,
		})
		return
	}

	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "get auth from context")
		return
	}

	if err := s.service.RatePerformance(r.Context(), userID, req.PerformanceID, req.Score, req.Comment, req.Gif); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) updatePerformance(w http.ResponseWriter, r *http.Request) {
	var req request.UpdatePerformanceRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{
			Err:  err.Error(),
			Code: ValidatationCode,
		})
		return
	}

	if err := s.service.UpdatePerformance(r.Context(), req.PerformanceID, req.Qualified); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getScoresFiltered(w http.ResponseWriter, r *http.Request) {
	var req request.GetScoresFilteredRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{
			Err:  err.Error(),
			Code: ValidatationCode,
		})
		return
	}

	scores, err := s.service.GetScoresFiltered(r.Context(), domain.Filters{
		UserID: req.UserID,
		CountryID: req.CountryID,
		ContestYear: req.ContestYear,
		Sort: req.Sort,
	})
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	EncodeJSONResponse(w, http.StatusOK, scores)
}

func (s *Server) getCountries(w http.ResponseWriter, r *http.Request) {
	res, err := s.service.GetCountries(r.Context())
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	EncodeJSONResponse(w, http.StatusOK, res)
}