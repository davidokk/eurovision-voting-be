package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(
				w,
				"missing authorization header",
				http.StatusUnauthorized,
			)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := s.jwt.ParseToken(tokenString)
		if err != nil {
			fmt.Println(err)
			http.Error(
				w,
				"invalid token",
				http.StatusUnauthorized,
			)
			return
		}

		ctx := setUserID(r.Context(), claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func (s *Server) VerifiedEmailMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := s.jwt.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		if !claims.EmailVerified {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(ApiError{
				Err:  "email not verified",
				Code: EmailNotVerifiedCode,
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

const userIDCtxKey = "user_id"

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDCtxKey, userID)
}

func getUserFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDCtxKey).(uuid.UUID)
	return id, ok 
}

func (s *Server) CheckRoleMw(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := s.jwt.ParseToken(tokenString)
			if err != nil {
				http.Error(
					w,
					"invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			user, err := s.service.GetUserByUsername(r.Context(), claims.Username)
			if err != nil {
				http.Error(w, fmt.Sprintf("cannot get user: %s", err.Error()), http.StatusInternalServerError)
				return 
			}
			
			if user.Role == nil || *user.Role != role {
				http.Error(w, fmt.Sprintf("role %s required", role), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}