package storage

type Config struct {
	URL        string `mapstructure:"url" valid:"required"`
	MinConns   int32  `mapstructure:"min_connections"`
	MaxConns   int32  `mapstructure:"max_connections" valid:"required"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
