package server

import (
	"errors"
	"eurovision-voting/internal/service"
	"net/http"
)

func mapAuthError(err error) (int, string) {
	if status, code, ok := mapTelegramAuthError(err); ok {
		return status, code
	}
	if status, code, ok := mapPasswordAuthError(err); ok {
		return status, code
	}
	switch {
	case errors.Is(err, service.ErrUsernameTaken):
		return http.StatusConflict, UsernameTakenCode
	case errors.Is(err, service.ErrInvalidCode):
		return http.StatusBadRequest, InvalidCodeCode
	default:
		return http.StatusInternalServerError, UnknownCode
	}
}

func authErrorMessage(err error, code string) string {
	switch code {
	case InvalidCodeCode:
		return "Неверный или просроченный код. Проверьте цифры или запросите новый код."
	case UsernameTakenCode:
		return "Это имя пользователя уже занято."
	case UserNotExistsCode:
		return "Пользователь не найден. Проверьте логин или зарегистрируйтесь."
	case WrongPasswordCode:
		return "Неверный пароль."
	case TelegramNotConfiguredCode:
		return "Вход через Telegram временно недоступен. Попробуйте позже."
	case TelegramRateLimitCode:
		return "Слишком много кодов в Telegram. Подождите около часа."
	case TelegramSessionInvalidCode:
		return "Сессия устарела. Запросите новую ссылку на сайте."
	case TelegramNotConnectedCode:
		return "Сначала откройте бота по ссылке и дождитесь кода в Telegram."
	case TelegramAlreadyLinkedCode:
		return "Этот Telegram уже привязан к другому аккаунту."
	case TelegramAccountNotFoundCode:
		return "Аккаунт с этим Telegram не найден. Зарегистрируйтесь."
	case SignupClosedCode:
		return "Регистрация закрыта :( Возвращайся в следующем году!"
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
