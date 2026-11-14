package main

import (
	"context"
	"fmt"
	jwtSvc "eurovision-voting/internal/jwt"
	"eurovision-voting/internal/media"
	"eurovision-voting/internal/server"
	"eurovision-voting/internal/service"
	"eurovision-voting/internal/storage"
	"eurovision-voting/internal/telegram"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultAddrPublic = ":8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sCh := listenToSignal()
	go func(sCh <-chan os.Signal, cancel func()) {
		<-sCh
		cancel()
	}(sCh, cancel)

	jwt := jwtSvc.NewJWTService(jwtSvc.Config{
		Secret:     os.Getenv("JWT_SECRET"),
		Expiration: 5 * 24 * time.Hour,
	})

	storage, err := storage.New(ctx, storage.Config{
		URL:      os.Getenv("POSTGRES_URL"),
		Username: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		MaxConns: 5,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("create storage")
	}

	s3, err := media.NewS3FromEnv(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("create S3")
	}

	signupAllowed := os.Getenv("SIGNUP_ALLOWED") != "0"
	svc := service.New(storage, jwt, s3, os.Getenv("TELEGRAM_BOT_USERNAME"), signupAllowed)

	srv := server.New(svc, jwt)

	var bot *telegram.Bot
	webhookBase := ""
	webhookSecret := ""

	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		bot, err = telegram.NewBot(token, svc)
		if err != nil {
			log.Error().Err(err).Msg("telegram bot disabled")
			bot = nil
		} else {
			svc.SetTelegramBotUsername(bot.Username())
			webhookBase = strings.TrimRight(strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_BASE_URL")), "/")
			webhookSecret = strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_SECRET"))
			if webhookBase != "" {
				srv.SetTelegramWebhook(bot.WebhookHandler())
			}
		}
	} else {
		log.Warn().Msg("TELEGRAM_BOT_TOKEN not set — Telegram auth disabled")
	}

	wg := sync.WaitGroup{}

	wg.Go(func() {
		log.Info().Msgf("Starting public server at: %s", defaultAddrPublic)
		if err := srv.ServePublic(ctx, defaultAddrPublic); err != nil {
			log.Error().Err(err).Msg("Could not serve public endpoints.")
		}
	})

	if bot != nil {
		if webhookBase != "" {
			webhookURL := fmt.Sprintf("%s/v1/telegram/webhook", webhookBase)
			if err := bot.SetWebhook(webhookURL, webhookSecret); err != nil {
				log.Fatal().Err(err).Str("url", webhookURL).Msg("telegram setWebhook")
			}
		} else {
			log.Warn().Msg("TELEGRAM_WEBHOOK_BASE_URL not set — falling back to long polling")
			wg.Go(func() {
				bot.RunPolling(ctx)
			})
		}
	}

	wg.Wait()
}

func listenToSignal() <-chan os.Signal {
	const ChanBuffer = 2
	sig := make(chan os.Signal, ChanBuffer)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return sig
}
