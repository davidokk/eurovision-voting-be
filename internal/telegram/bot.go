package telegram

import (
	"context"
	"errors"
	"eurovision-voting/internal/service"
	"fmt"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	svc         *service.Service
	secretToken string
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

// SetWebhook registers the HTTPS webhook with Telegram and stops long-polling.
// secret_token is sent via MakeRequest because telegram-bot-api v5.5.1's
// WebhookConfig does not expose the field yet.
func (b *Bot) SetWebhook(webhookURL, secretToken string) error {
	b.secretToken = strings.TrimSpace(secretToken)
	webhookURL = strings.TrimSpace(webhookURL)
	if webhookURL == "" {
		return fmt.Errorf("telegram webhook url is empty")
	}

	params := tgbotapi.Params{
		"url":                  webhookURL,
		"max_connections":      "40",
		"drop_pending_updates": "true",
	}
	if b.secretToken != "" {
		params["secret_token"] = b.secretToken
	}

	resp, err := b.api.MakeRequest("setWebhook", params)
	if err != nil {
		return fmt.Errorf("telegram setWebhook: %w", err)
	}
	if resp != nil && !resp.Ok {
		return fmt.Errorf("telegram setWebhook not ok: %s", resp.Description)
	}

	info, err := b.api.GetWebhookInfo()
	if err != nil {
		log.Warn().Err(err).Msg("telegram getWebhookInfo")
	} else {
		log.Info().
			Str("url", info.URL).
			Int("pending", info.PendingUpdateCount).
			Msg("telegram webhook registered")
	}
	return nil
}

// DeleteWebhook removes the webhook (useful for local polling fallback).
func (b *Bot) DeleteWebhook() error {
	_, err := b.api.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
	return err
}

// WebhookHandler processes Telegram Update POSTs.
func (b *Bot) WebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if b.secretToken != "" {
			got := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if got != b.secretToken {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		update, err := b.api.HandleUpdate(r)
		if err != nil {
			log.Error().Err(err).Msg("telegram handle update")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		b.handleUpdate(r.Context(), *update)
		w.WriteHeader(http.StatusOK)
	}
}

// RunPolling is a local-dev fallback when webhook URL is not configured.
func (b *Bot) RunPolling(ctx context.Context) {
	if err := b.DeleteWebhook(); err != nil {
		log.Warn().Err(err).Msg("telegram deleteWebhook before polling")
	}

	log.Info().Str("username", b.api.Self.UserName).Msg("telegram bot polling started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
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

	text := fmt.Sprintf("Твой код для входа:\n<code>%s</code>", code)
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
