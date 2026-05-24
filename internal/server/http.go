package server

import (
	"context"
	"encoding/json"
	"eurovision-voting/internal/jwt"
	"eurovision-voting/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	service      *service.Service
	jwt          *jwt.Service
	publicRouter *chi.Mux
}

func New(service *service.Service, jwt *jwt.Service) *Server {
	srv := &Server{
		service: service,
		jwt:     jwt,

		publicRouter: chi.NewRouter(),
	}

	srv.publicRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"https://eurovision-voting-fe.onrender.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	publicMws := []func(http.Handler) http.Handler{}

	srv.registerPublicRoutes(publicMws...)

	return srv
}

func (s *Server) registerPublicRoutes(mws ...func(http.Handler) http.Handler) {
	s.publicRouter.Use(mws...)

	s.publicRouter.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/telegram/signin/start", s.telegramSigninStartHandler)
			r.Post("/telegram/signin/confirm", s.telegramSigninConfirmHandler)
			r.Get("/telegram/session", s.telegramSessionStatusHandler)
			r.With(s.AuthMiddleware).Get("/validate", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/user", func(r chi.Router) {
			r.With(s.AuthMiddleware).Get("/me", s.getMe)
			r.With(s.AuthMiddleware).Patch("/username", s.changeUsernameHandler)
			r.With(s.AuthMiddleware, s.VerifiedAccountMiddleware).Delete("/avatar", s.deleteAvatar)
			r.Get("/public", s.getUserPublic)
			r.Get("/list", s.listUsersHandler)
		})

		r.With(s.AuthMiddleware, s.VerifiedAccountMiddleware).Post("/media/upload", s.uploadMedia)

		r.Route("/contest", func(r chi.Router) {
			r.Get("/", s.getContestsByYearHandler)
			r.Get("/{id}", s.getContestViewHandler)
		})

		r.Get("/scores", s.getScoresFiltered)
		r.Get("/countries", s.getCountries)

		r.Route("/proxy", func(r chi.Router) {
			r.Get("/youtube/search", s.proxyYouTubeSearch)
			r.Get("/giphy/search", s.proxyGiphySearch)
		})

		r.With(s.AuthMiddleware, s.VerifiedAccountMiddleware).Post("/performance/{id}/rate", s.ratePerformance)

		r.Get("/ws", s.serveWS())

		r.Route("/message", func(r chi.Router) {
			r.With(s.AuthMiddleware, s.VerifiedAccountMiddleware).Post("/send", s.sendMessage)
			r.Get("/", s.getMessages)
		})
	})

	s.publicRouter.Route("/admin", func(r chi.Router) {
		r.Use(s.CheckRoleMw("ADMIN"))
		r.Post("/contest", s.adminCreateContest)
		r.Put("/contest/{id}", s.adminUpdateContest)
		r.Post("/contest/{id}/performance", s.adminCreatePerformance)
		r.Put("/contest/{id}/places", s.adminUpdatePlaces)
		r.Put("/performance/{id}", s.updatePerformance)
	})
}

func (s *Server) ServePublic(ctx context.Context, addr string) error {
	return s.serve(ctx, addr, s.publicRouter)
}

func (s *Server) serve(ctx context.Context, addr string, h http.Handler) error {
	server := http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	errCh := make(chan error)
	go func(errCh chan<- error) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}(errCh)

	select {
	case <-ctx.Done():
		log.Info().Msg("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func EncodeJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}
