package server

import (
	"eurovision-voting/internal/server/request"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (s *Server) adminCreateContest(w http.ResponseWriter, r *http.Request) {
	var req request.AdminCreateContestRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	starts, err := time.Parse(time.RFC3339, req.Starts)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid starts", Code: ValidatationCode})
		return
	}
	ends, err := time.Parse(time.RFC3339, req.Ends)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid ends", Code: ValidatationCode})
		return
	}

	c, err := s.service.AdminCreateContest(r.Context(), req.Year, req.Type, starts, ends)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	EncodeJSONResponse(w, http.StatusCreated, c)
}

func (s *Server) adminUpdateContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid contest id", Code: ValidatationCode})
		return
	}

	var req request.AdminCreateContestRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	starts, err := time.Parse(time.RFC3339, req.Starts)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid starts", Code: ValidatationCode})
		return
	}
	ends, err := time.Parse(time.RFC3339, req.Ends)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid ends", Code: ValidatationCode})
		return
	}

	c, err := s.service.AdminUpdateContest(r.Context(), contestID, req.Year, req.Type, starts, ends)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, c)
}

func (s *Server) adminCreatePerformance(w http.ResponseWriter, r *http.Request) {
	contestID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid contest id", Code: ValidatationCode})
		return
	}

	var req request.AdminCreatePerformanceRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	countryID, _ := uuid.Parse(req.CountryID)
	perfID, err := s.service.AdminCreatePerformance(
		r.Context(), contestID, countryID, req.Number,
		req.Artist, req.Song, req.YoutubeLink,
	)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	EncodeJSONResponse(w, http.StatusCreated, map[string]string{"id": perfID.String()})
}

func (s *Server) adminUpdatePlaces(w http.ResponseWriter, r *http.Request) {
	contestID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid contest id", Code: ValidatationCode})
		return
	}

	var req request.AdminUpdatePlacesRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	ranked := make([]uuid.UUID, 0, len(req.Ranked))
	for _, idStr := range req.Ranked {
		id, err := uuid.Parse(idStr)
		if err != nil {
			EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid performance id in ranked", Code: ValidatationCode})
			return
		}
		ranked = append(ranked, id)
	}

	if err := s.service.AdminUpdatePlaces(r.Context(), contestID, ranked); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
