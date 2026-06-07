package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	IsDebugLoggerConfig bool `env:"IS_DEBUG_LOGGER_CONFIG" envDefault:"true"`
}

func NewConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}
	return cfg
}
