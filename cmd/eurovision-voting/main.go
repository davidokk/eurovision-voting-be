package main

import (
	"context"
	jwtSvc "eurovision-voting/internal/jwt"
	"eurovision-voting/internal/media"
	"eurovision-voting/internal/server"
	"eurovision-voting/internal/service"
	"eurovision-voting/internal/storage"
	"os"
	"os/signal"
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
		Expiration: 24 * time.Hour,
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

	svc := service.New(storage, jwt, s3)

	srv := server.New(svc, jwt)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		log.Info().Msgf("Starting public server at: %s", defaultAddrPublic)
		if err := srv.ServePublic(ctx, defaultAddrPublic); err != nil {
			log.Error().Err(err).Msg("Could not serve public endpoints.")
		}
	})

	wg.Wait()
}

func listenToSignal() <-chan os.Signal {
	const ChanBuffer = 2
	sig := make(chan os.Signal, ChanBuffer)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return sig
}
