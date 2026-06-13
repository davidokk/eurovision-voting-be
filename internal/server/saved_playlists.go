package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/service"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type savedPlaylistBody struct {
	Name    string                          `json:"name"`
	Entries []domain.GamePlaylistEntryInput `json:"entries"`
}

func (s *Server) listSavedPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := s.service.ListSavedPlaylists(r.Context(), userID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, items)
}

func (s *Server) getSavedPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid id"})
		return
	}
	item, err := s.service.GetSavedPlaylist(r.Context(), userID, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, item)
}

func (s *Server) createSavedPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req savedPlaylistBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid body"})
		return
	}
	item, err := s.service.CreateSavedPlaylist(r.Context(), userID, req.Name, req.Entries)
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusCreated, item)
}

func (s *Server) updateSavedPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid id"})
		return
	}
	var req savedPlaylistBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid body"})
		return
	}
	item, err := s.service.UpdateSavedPlaylist(r.Context(), userID, id, req.Name, req.Entries)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, item)
}

func (s *Server) deleteSavedPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid id"})
		return
	}
	if err := s.service.DeleteSavedPlaylist(r.Context(), userID, id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
