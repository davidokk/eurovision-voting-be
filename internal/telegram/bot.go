package telegram

import (
	"context"
	"errors"
	"eurovision-voting/internal/service"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	api *tgbotapi.BotAPI
	svc *service.Service
}

func NewBot(token string, svc *service.Service) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("telegram bot: %w", err)
	}
	api.Debug = false
	return &Bot{api: api, svc: svc}, nil
}

func (b *Bot) Username() string {
	return b.api.Self.UserName
}

func (b *Bot) Run(ctx context.Context) {
	log.Info().Str("username", b.api.Self.UserName).Msg("telegram bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			b.handleUpdate(ctx, update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	msg := update.Message
	switch msg.Command() {
	case "start":
		b.handleStart(ctx, msg)
	default:
		b.reply(msg.Chat.ID, "Используйте ссылку с сайта Eurovision Voting, чтобы получить код входа.")
	}
}

func (b *Bot) handleStart(ctx context.Context, msg *tgbotapi.Message) {
	linkToken := strings.TrimSpace(msg.CommandArguments())
	if linkToken == "" {
		b.reply(msg.Chat.ID, "Привет! Откройте ссылку «Войти через Telegram» на сайте eurovision-voting — бот пришлёт код для входа.")
		return
	}

	if msg.From == nil {
		b.reply(msg.Chat.ID, "Не удалось определить ваш Telegram-профиль.")
		return
	}

	tgUsername := ""
	if msg.From.UserName != "" {
		tgUsername = msg.From.UserName
	}

	code, err := b.svc.BotDeliverCode(ctx, linkToken, msg.From.ID, msg.Chat.ID, tgUsername)
	if err != nil {
		b.reply(msg.Chat.ID, botErrorMessage(err))
		return
	}

	text := fmt.Sprintf("Ваш код для входа -\n\n<code>%s</code>", code)
	out := tgbotapi.NewMessage(msg.Chat.ID, text)
	out.ParseMode = "HTML"
	if _, err := b.api.Send(out); err != nil {
		log.Error().Err(err).Msg("telegram send code")
	}
}

func (b *Bot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Error().Err(err).Msg("telegram reply")
	}
}

func botErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrTelegramSessionInvalid):
		return "Ссылка устарела или уже использована. Вернитесь на сайт и запросите новую."
	case errors.Is(err, service.ErrTelegramRateLimit):
		return "Слишком много запросов. Подождите около часа и попробуйте снова."
	case errors.Is(err, service.ErrTelegramAlreadyLinked):
		return "Этот Telegram уже привязан к другому аккаунту."
	case errors.Is(err, service.ErrSignupClosed):
		return "Регистрация закрыта. Возвращайся в следующем году!"
	default:
		log.Error().Err(err).Msg("telegram bot deliver code")
		return "Не удалось выдать код. Попробуйте позже или запросите новую ссылку на сайте."
	}
}
