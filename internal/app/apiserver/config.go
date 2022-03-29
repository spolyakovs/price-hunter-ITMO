package apiserver

type Config struct {
	BindAddr      string `toml:"bind_addr"`
	LogLevel      string `toml:"log_level"`
	DatabaseURL   string `toml:"database_url"`
	AccessSecret  string `toml:"access_secret"`
	RefreshSecret string `toml:"refresh_secret"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
