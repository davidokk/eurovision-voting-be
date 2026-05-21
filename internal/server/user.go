package server

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := s.service.GetUserPublic(r.Context(), userID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	EncodeJSONResponse(w, http.StatusOK, u)
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
