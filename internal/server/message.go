package server

import (
	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/server/request"
	"net/http"
)

func (s *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req request.SendMessageRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	msg := &domain.Message{
		ContestID:       req.ContestID,
		Message:         req.Message,
		UserID:          userID,
		ContentType:     req.ContentType,
		MediaURL:        req.MediaURL,
		MediaDurationMs: req.MediaDurationMs,
		ReplyToID:       req.ReplyToID,
	}
	if err := s.service.SendMessage(r.Context(), msg); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	EncodeJSONResponse(w, http.StatusOK, msg)
}
