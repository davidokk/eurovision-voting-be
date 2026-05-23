package server

import (
	"errors"
	"eurovision-voting/internal/server/request"
	"eurovision-voting/internal/server/response"
	"eurovision-voting/internal/service"
	"net/http"
	"os"
)

func (s *Server) signupStartHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SIGNUP_ALLOWED") == "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var req request.SignupStartRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	if err := s.service.StartSignup(r.Context(), req.Email, req.Username, req.Password); err != nil {
		writeAuthError(w, err)
		return
	}

	EncodeJSONResponse(w, http.StatusOK, map[string]string{"status": "code_sent"})
}

func (s *Server) signupConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SIGNUP_ALLOWED") == "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var req request.SignupConfirmRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	token, err := s.service.ConfirmSignup(r.Context(), req.Email, req.Code)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	EncodeJSONResponse(w, http.StatusOK, response.AuthResponse{Token: token})
}

func (s *Server) signinHandler(w http.ResponseWriter, r *http.Request) {
	var req request.SignInRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	token, me, err := s.service.SignIn(r.Context(), req.Identifier, req.Password)
	if err != nil {
		apiErr := ApiError{Err: err.Error()}
		statusCode := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrUserNotExists):
			apiErr.Code = UserNotExistsCode
			statusCode = http.StatusNotFound
		case errors.Is(err, service.ErrWrongPassword):
			apiErr.Code = WrongPasswordCode
			statusCode = http.StatusForbidden
		default:
			apiErr.Code = UnknownCode
		}
		EncodeJSONResponse(w, statusCode, apiErr)
		return
	}

	EncodeJSONResponse(w, http.StatusOK, response.AuthResponse{
		Token: token,
		User:  me,
	})
}

func (s *Server) emailBindRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req request.EmailRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	if err := s.service.RequestEmailBind(r.Context(), userID, req.Email); err != nil {
		writeAuthError(w, err)
		return
	}

	EncodeJSONResponse(w, http.StatusOK, map[string]string{"status": "code_sent"})
}

func (s *Server) emailBindConfirmHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserFromContext(r.Context())
	if !ok {
		EncodeJSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req request.EmailConfirmRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	token, err := s.service.ConfirmEmailBind(r.Context(), userID, req.Email, req.Code)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	EncodeJSONResponse(w, http.StatusOK, response.AuthResponse{Token: token})
}

func (s *Server) passwordForgotHandler(w http.ResponseWriter, r *http.Request) {
	var req request.EmailRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	if err := s.service.RequestPasswordReset(r.Context(), req.Email); err != nil {
		writeAuthError(w, err)
		return
	}
	EncodeJSONResponse(w, http.StatusOK, map[string]string{"status": "code_sent"})
}

func (s *Server) passwordResetHandler(w http.ResponseWriter, r *http.Request) {
	var req request.PasswordResetRequest
	if err := req.Bind(r); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}
	if err := req.Validate(); err != nil {
		EncodeJSONResponse(w, http.StatusBadRequest, ApiError{Err: err.Error(), Code: ValidatationCode})
		return
	}

	if err := s.service.ConfirmPasswordReset(r.Context(), req.Email, req.Code, req.Password); err != nil {
		writeAuthError(w, err)
		return
	}

	EncodeJSONResponse(w, http.StatusOK, map[string]string{"status": "password_updated"})
}

func mapAuthError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrEmailTaken):
		return http.StatusConflict, EmailTakenCode
	case errors.Is(err, service.ErrUsernameTaken):
		return http.StatusConflict, UsernameTakenCode
	case errors.Is(err, service.ErrInvalidCode):
		return http.StatusBadRequest, InvalidCodeCode
	case errors.Is(err, service.ErrEmailRateLimit):
		return http.StatusTooManyRequests, EmailRateLimitCode
	case errors.Is(err, service.ErrInvalidEmail):
		return http.StatusBadRequest, ValidatationCode
	case errors.Is(err, service.ErrEmailNotConfigured):
		return http.StatusServiceUnavailable, EmailNotConfiguredCode
	default:
		return http.StatusInternalServerError, UnknownCode
	}
}

func authErrorMessage(err error, code string) string {
	switch code {
	case InvalidCodeCode:
		return "Неверный или просроченный код. Проверьте цифры или запросите новый код."
	case EmailRateLimitCode:
		return "Слишком много писем на этот адрес. Можно отправить не больше 5 писем в час — попробуйте позже."
	case EmailTakenCode:
		return "Этот email уже зарегистрирован."
	case UsernameTakenCode:
		return "Это имя пользователя уже занято."
	case EmailNotConfiguredCode:
		return "Отправка писем временно недоступна. Попробуйте позже."
	default:
		return err.Error()
	}
}

func writeAuthError(w http.ResponseWriter, err error) {
	status, code := mapAuthError(err)
	EncodeJSONResponse(w, status, ApiError{
		Err:  authErrorMessage(err, code),
		Code: code,
	})
}
