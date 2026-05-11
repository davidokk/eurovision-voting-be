package server

import (
	"context"
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

const userIDCtxKey = "user_id"

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDCtxKey, userID)
}

func getUserFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDCtxKey).(uuid.UUID)
	return id, ok 
}