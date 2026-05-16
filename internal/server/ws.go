package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // dev only
	},
}

func (s *Server) serveWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("upgrade")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("upgrade conn")
			return
		}

		userID := uuid.New()
		if tok := r.URL.Query().Get("token"); tok != "" {
			if claims, err := s.jwt.ParseToken(tok); err == nil {
				userID = claims.UserID
			}
		}

		log.Info().Str("user", userID.String()).Msg("serve conn")
		s.service.ServeConn(uuid.New(), conn)
	}

}
