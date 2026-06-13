package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"eurovision-voting/internal/domain"
	"eurovision-voting/internal/service"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (s *Server) getGameCatalog(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.GetGameCatalog(r.Context())
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	if items == nil {
		items = []domain.GameCatalogItem{}
	}
	EncodeJSONResponse(w, http.StatusOK, items)
}

func (s *Server) createGameRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	room, err := s.service.CreateGameRoom(r.Context(), userID)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusCreated, room)
}

func (s *Server) joinGameRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	code := chi.URLParam(r, "code")
	room, err := s.service.JoinGameRoom(r.Context(), code, userID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, room)
}

func (s *Server) getGameRoom(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	userID, authed := getUserFromContext(r.Context())
	var room *domain.GameRoomView
	var err error
	if authed {
		room, err = s.service.GetGameRoomForUser(code, userID)
	} else {
		room, err = s.service.GetGameRoom(code)
	}
	if err != nil {
		status := http.StatusNotFound
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, room)
}

type gameActionRequest struct {
	Action  string         `json:"action"`
	Payload map[string]any `json:"payload"`
}

func (s *Server) gameRoomAction(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	code := chi.URLParam(r, "code")

	var req gameActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid body"})
		return
	}
	if req.Payload == nil {
		req.Payload = map[string]any{}
	}

	room, err := s.service.HandleGameHostAction(r.Context(), code, userID, req.Action, req.Payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrGameForbidden) {
			status = http.StatusForbidden
		} else if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, room)
}

func (s *Server) gameBuzz(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	code := chi.URLParam(r, "code")
	room, err := s.service.GameBuzz(code, userID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, room)
}

type gameAnswerRequest struct {
	Answer string `json:"answer"`
}

func (s *Server) gameSubmitAnswer(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	code := chi.URLParam(r, "code")

	var req gameAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "invalid body"})
		return
	}

	room, err := s.service.GameSubmitAnswer(code, userID, req.Answer)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrGameNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrGameForbidden) {
			status = http.StatusForbidden
		}
		EncodeJSONResponse(w, status, ApiError{Err: err.Error()})
		return
	}
	EncodeJSONResponse(w, http.StatusOK, room)
}

func (s *Server) serveGameWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("room")))
		if code == "" {
			_ = conn.WriteJSON(map[string]string{"type": "game.error", "message": "room required"})
			_ = conn.Close()
			return
		}

		userID, username, ok := s.parseWSUser(r)
		if !ok {
			_ = conn.WriteJSON(map[string]string{"type": "game.error", "message": "auth required"})
			_ = conn.Close()
			return
		}

		s.service.ServeGameConn(code, userID, username, conn)
	}
}

func (s *Server) parseWSUser(r *http.Request) (userID uuid.UUID, username string, ok bool) {
	tok := r.URL.Query().Get("token")
	if tok == "" {
		return uuid.Nil, "", false
	}
	claims, err := s.jwt.ParseToken(tok)
	if err != nil {
		return uuid.Nil, "", false
	}
	return claims.UserID, claims.Username, true
}
