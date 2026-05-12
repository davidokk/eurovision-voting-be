package server

import (
	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/server/request"
	"net/http"
)

func (s *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "get auth from context")
		return
	}

	var req request.SendMessageRequest
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

	s.service.SendMessage(r.Context(), &domain.Message{
		ContestID: req.ContestID,
		Message:   req.Message,
		UserID:    userID,
	})

	w.WriteHeader(http.StatusOK)
}
