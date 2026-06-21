package server

import (
	"net/http"

	"eurovision-voting/internal/domain"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (s *Server) listFavoritePerformanceIDs(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ids, err := s.service.ListFavoritePerformanceIDs(r.Context(), userID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	if ids == nil {
		ids = []string{}
	}
	EncodeJSONResponse(w, http.StatusOK, ids)
}

func (s *Server) listFavoritePerformances(w http.ResponseWriter, r *http.Request) {
	targetID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid user_id"})
		return
	}
	items, err := s.service.ListFavoritePerformances(r.Context(), targetID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	if items == nil {
		items = []domain.FavoritePerformance{}
	}
	EncodeJSONResponse(w, http.StatusOK, items)
}

func (s *Server) addPerformanceFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	performanceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid performance id"})
		return
	}
	if err := s.service.AddPerformanceFavorite(r.Context(), userID, performanceID); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) removePerformanceFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	performanceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid performance id"})
		return
	}
	if err := s.service.RemovePerformanceFavorite(r.Context(), userID, performanceID); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
}
