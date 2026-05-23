package server

import (
	"errors"
	"eurovision-voting/internal/server/request"
	"eurovision-voting/internal/server/response"
	"eurovision-voting/internal/service"
	"net/http"
	"os"
)

func (s *Server) telegramSigninStartHandler(w http.ResponseWriter, r *http.Request) {
	res, err := s.service.StartTelegramSignin(r.Context())
	if err != nil {
		writeAuthError(w, err)
		return
	}
	EncodeJSONResponse(w, http.StatusOK, res)
}

func (s *Server) telegramSigninConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var req request.TelegramConfirmRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	signupAllowed := os.Getenv("SIGNUP_ALLOWED") != "0"
	token, me, err := s.service.ConfirmTelegramSignin(r.Context(), req.LinkToken, req.Code, signupAllowed)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	EncodeJSONResponse(w, http.StatusOK, response.AuthResponse{Token: token, User: me})
}

func (s *Server) telegramSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	linkToken := r.URL.Query().Get("link_token")
	if linkToken == "" {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: "link_token required", Code: ValidatationCode})
		return
	}

	st, err := s.service.GetTelegramSessionStatus(r.Context(), linkToken)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	EncodeJSONResponse(w, http.StatusOK, st)
}

func mapTelegramAuthError(err error) (int, string, bool) {
	switch {
	case errors.Is(err, service.ErrSignupClosed):
		return http.StatusForbidden, SignupClosedCode, true
	case errors.Is(err, service.ErrTelegramNotConfigured):
		return http.StatusServiceUnavailable, TelegramNotConfiguredCode, true
	case errors.Is(err, service.ErrTelegramRateLimit):
		return http.StatusTooManyRequests, TelegramRateLimitCode, true
	case errors.Is(err, service.ErrTelegramSessionInvalid):
		return http.StatusBadRequest, TelegramSessionInvalidCode, true
	case errors.Is(err, service.ErrTelegramNotConnected):
		return http.StatusBadRequest, TelegramNotConnectedCode, true
	case errors.Is(err, service.ErrTelegramAlreadyLinked):
		return http.StatusConflict, TelegramAlreadyLinkedCode, true
	default:
		return 0, "", false
	}
}
