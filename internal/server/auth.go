package server

import (
	"errors"
	"eurovision-voting/internal/server/request"
	"eurovision-voting/internal/server/response"
	"eurovision-voting/internal/service"
	"net/http"
)

func (s *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	var req request.SingUpRequest
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

	token, err := s.service.SignUp(r.Context(), req.Username, req.Password)
	if err != nil {
		EncodeJSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	EncodeJSONResponse(w, http.StatusOK, response.SignUpResponse{
		Token: token,
	})
}

func (s *Server) signinHandler(w http.ResponseWriter, r *http.Request) {
	var req request.SingUpRequest
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

	token, err := s.service.SignIn(r.Context(), req.Username, req.Password)
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

	EncodeJSONResponse(w, http.StatusOK, response.SignUpResponse{
		Token: token,
	})
}
