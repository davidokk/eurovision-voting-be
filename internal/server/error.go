package server

const (
	ValidatationCode  = "VALIDATE"
	UserNotExistsCode = "USER_NOT_EXISTS"
	WrongPasswordCode = "WRONG_PASSWORD"
	UnknownCode = "UNKNOWN"
)

type ApiError struct {
	Err  string `json:"message"`
	Code string `json:"code"`
}
