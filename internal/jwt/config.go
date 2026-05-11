package jwt

import "time"

type Config struct {
	Secret     string
	Expiration time.Duration
}
