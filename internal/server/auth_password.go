package server

import (
	"errors"
	"net/http"
	"os"

	"eurovision-voting/internal/server/request"
	"eurovision-voting/internal/server/response"
	"eurovision-voting/internal/service"
)

func (s *Server) signinHandler(w http.ResponseWriter, r *http.Request) {
	var req request.CredentialsRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	token, me, err := s.service.SignInWithPassword(r.Context(), req.Username, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	EncodeJSONResponse(w, http.StatusOK, response.AuthResponse{Token: token, User: me})
}

func (s *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SIGNUP_ALLOWED") == "0" {
		writeAuthError(w, service.ErrSignupClosed)
		return
	}

	var req request.CredentialsRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	token, me, err := s.service.SignUpWithPassword(r.Context(), req.Username, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	EncodeJSONResponse(w, http.StatusCreated, response.AuthResponse{Token: token, User: me})
}

func mapPasswordAuthError(err error) (int, string, bool) {
	switch {
	case errors.Is(err, service.ErrUserNotExists):
		return http.StatusUnauthorized, UserNotExistsCode, true
	case errors.Is(err, service.ErrWrongPassword):
		return http.StatusUnauthorized, WrongPasswordCode, true
	case errors.Is(err, service.ErrSignupClosed):
		return http.StatusForbidden, SignupClosedCode, true
	default:
		return 0, "", false
	}
}
