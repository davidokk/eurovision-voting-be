package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret          string
	tokenExpiration time.Duration
}

func NewJWTService(cfg Config) *Service {
	return &Service{
		secret:          cfg.Secret,
		tokenExpiration: cfg.Expiration,
	}
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Verified bool `json:"verified"`
	jwt.RegisteredClaims
}

func (c *Claims) IsVerified() bool {
	return c.Verified
}

func (s *Service) GenerateToken(userID uuid.UUID, username string, verified bool) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Verified: verified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(s.tokenExpiration),
			),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString([]byte(s.secret))

}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil

}
