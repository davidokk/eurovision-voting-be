package server

const (
	ValidatationCode            = "VALIDATE"
	UserNotExistsCode           = "USER_NOT_EXISTS"
	WrongPasswordCode           = "WRONG_PASSWORD"
	AccountNotVerifiedCode      = "ACCOUNT_NOT_VERIFIED"
	UsernameTakenCode           = "USERNAME_TAKEN"
	InvalidCodeCode             = "INVALID_CODE"
	TelegramNotConfiguredCode   = "TELEGRAM_NOT_CONFIGURED"
	TelegramRateLimitCode       = "TELEGRAM_RATE_LIMIT"
	TelegramSessionInvalidCode  = "TELEGRAM_SESSION_INVALID"
	TelegramNotConnectedCode    = "TELEGRAM_NOT_CONNECTED"
	TelegramAlreadyLinkedCode   = "TELEGRAM_ALREADY_LINKED"
	TelegramAccountNotFoundCode = "TELEGRAM_ACCOUNT_NOT_FOUND"
	SignupClosedCode            = "SIGNUP_CLOSED"
	UnknownCode                 = "UNKNOWN"
)

type ApiError struct {
	Err  string `json:"message"`
	Code string `json:"code"`
}
