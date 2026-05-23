package server

const (
	ValidatationCode       = "VALIDATE"
	UserNotExistsCode      = "USER_NOT_EXISTS"
	WrongPasswordCode      = "WRONG_PASSWORD"
	EmailNotVerifiedCode   = "EMAIL_NOT_VERIFIED"
	EmailTakenCode         = "EMAIL_TAKEN"
	UsernameTakenCode      = "USERNAME_TAKEN"
	InvalidCodeCode        = "INVALID_CODE"
	EmailNotConfiguredCode = "EMAIL_NOT_CONFIGURED"
	EmailRateLimitCode     = "EMAIL_RATE_LIMIT"
	UnknownCode            = "UNKNOWN"
)

type ApiError struct {
	Err  string `json:"message"`
	Code string `json:"code"`
}
