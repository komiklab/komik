package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	IsDebugLoggerConfig bool   `env:"IS_DEBUG_LOGGER_CONFIG" envDefault:"true"`
	PostgresDSN         string `env:"POSTGRES_DSN" envDefault:"postgres://komik:komik@localhost:5433/komik?sslmode=disable"`
	CORSSupport         string `env:"CORS_SUPPORT" envDefault:"http://localhost:3000"`
}

func NewConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}
	return cfg
}
