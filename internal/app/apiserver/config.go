package apiserver

type Config struct {
	BindAddr         string `toml:"BIND_ADDR"`
	LogLevel         string `toml:"LOG_LEVEL"`
	DatabaseHost     string `toml:"DATABASE_HOST"`
	DatabaseDBName   string `toml:"DATABASE_DB"`
	DatabaseUser     string `toml:"DATABASE_USER"`
	DatabasePassword string `toml:"DATABASE_PASSWORD"`
	DatabaseSSLMode  string `toml:"DATABASE_SSLMODE"`
	RedisAddr        string `toml:"REDIS_ADDR"`
	TokenSecret      string `toml:"TOKEN_SECRET"`
	SteamAPIKey      string `toml:"STEAM_API_KEY"`
}

// TODO: think about changing default config
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
