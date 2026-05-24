package server

import (
	"net/http"
	"strings"

	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/server/request"
	"eurovision-voting/internal/server/response"

	"github.com/google/uuid"
)

func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := s.service.GetUserMe(r.Context(), userID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	EncodeJSONResponse(w, http.StatusOK, u)
}

func (s *Server) changeUsernameHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req request.ChangeUsernameRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	token, err := s.service.ChangeUsername(r.Context(), userID, req.Username)
	if err != nil {
		status, code := mapAuthError(err)
		EncodeJSONResponse(w, status, ApiError{Err: err.Error(), Code: code})
		return
	}

	EncodeJSONResponse(w, http.StatusOK, response.AuthResponse{Token: token})
}

func (s *Server) deleteAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := s.service.DeleteUserAvatar(r.Context(), userID); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	EncodeJSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	var excludeID *uuid.UUID
	if ex := strings.TrimSpace(r.URL.Query().Get("exclude")); ex != "" {
		if id, err := uuid.Parse(ex); err == nil {
			excludeID = &id
		}
	}
	users, err := s.service.ListUsers(r.Context(), excludeID, 1000)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []domain.UserPublic{}
	}
	EncodeJSONResponse(w, http.StatusOK, users)
}

func (s *Server) getUserPublic(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid id"})
		return
	}
	u, err := s.service.GetUserPublic(r.Context(), id)
	if err != nil {
		EncodeJSONResponse(w, http.StatusNotFound, err.Error())
		return
	}
	EncodeJSONResponse(w, http.StatusOK, u)
}
