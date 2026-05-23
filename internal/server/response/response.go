package response

import "eurovision-voting/internal/domain"

type SignUpResponse struct {
	Token string `json:"token"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	Token string          `json:"token"`
	User  *domain.UserMe  `json:"user,omitempty"`
}

