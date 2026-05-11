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
		AllowedOrigins:   []string{"http://localhost:5173", "https://eurovision-voting-fe.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
			r.Post("/signup", s.signupHandler)
			r.Post("/signin", s.signinHandler)
		})

		r.Route("/contest", func(r chi.Router) {
			r.Get("/", s.getContestsByYearHandler)
			r.Get("/{id}", s.getContestViewHandler)
		})

		r.With(s.AuthMiddleware).Post("/performance/{id}/rate", s.ratePerformance)
	})
}

func (s *Server) ServePublic(ctx context.Context, addr string) error {
	return s.serve(ctx, addr, s.publicRouter)
}

func (s *Server) serve(ctx context.Context, addr string, h http.Handler) error {
	server := http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func(errCh chan<- error) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}(errCh)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info().Msg("Server is interrupted. Exiting...")
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Close()
	}

	return nil
}

func EncodeJSONResponse[T any](w http.ResponseWriter, code int, data T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(data)
}
