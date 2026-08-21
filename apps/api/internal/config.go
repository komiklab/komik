package internal

import (
	"github.com/asaskevich/govalidator/v12"
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

type Config struct {
	IsDebugLoggerConfig bool   `valid:"-" env:"IS_DEBUG_LOGGER_CONFIG" envDefault:"true"`
	PostgresDSN         string `valid:"-" env:"POSTGRES_DSN" envDefault:"postgres://komik:komik@localhost:5433/komik?sslmode=disable"`
	CORSSupport         string `valid:"-" env:"CORS_SUPPORT" envDefault:"http://localhost:3000"`
	NatsURL             string `valid:"-" env:"NATS_URL" envDefault:"nats://localhost:4223"`
	RedisDSN            string `valid:"-" env:"REDIS_DSN" envDefault:"redis://localhost:6381"`
	PostLoginRedirect   string `valid:"-" env:"POST_LOGIN_REDIRECT" envDefault:"http://localhost:3000/"`
	TemporalConfig      TemporalConfig
	OauthConfig         OauthConfig
	SlackIntegration    SlackIntegration
	InngestConfig       InngestConfig
}

type InngestConfig struct {
	EventKey   string `valid:"required" env:"INNGEST_EVENT_KEY" envDefault:""`
	SigningKey string `valid:"required" env:"INNGEST_SIGNING_KEY" envDefault:""`
	AppID      string `env:"INNGEST_APP_ID" envDefault:"komik-app"`
	AppVersion string `env:"INNGEST_APP_VERSION" envDefault:""`
	Host       string `env:"INNGEST_HOST" envDefault:"localhost:8288"`
}

type TemporalConfig struct {
	Namespace   string `valid:"-" env:"TEMPORAL_NAMESPACE" envDefault:"default"`
	TemporalUrl string `valid:"-" env:"TEMPORAL_URL" envDefault:"localhost:7233"`
	TaskQueue   string `valid:"-" env:"TEMPORAL_TASK_QUEUE" envDefault:"komik-TaskQueue"`
}

type OauthConfig struct {
	ClientID     string `valid:"required" env:"CLIENT_ID"`
	ClientSecret string `valid:"required" env:"CLIENT_SECRET"`
	RedirectURL  string `valid:"-" env:"REDIRECT_URL"`
	Scopes       string `valid:"-" env:"SCOPES" envDefault:""`
	AuthURL      string `valid:"-" env:"AUTH_URL" envDefault:""`
	TokenURL     string `valid:"-" env:"TOKEN_URL" envDefault:""`
}

type SlackIntegration struct {
	BotToken      string `valid:"-" env:"SLACK_BOT_TOKEN" envDefault:""`
	SigningSecret string `valid:"-" env:"SLACK_SIGNING_SECRET" envDefault:""`
}

func NewConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to validate config")
	}

	return cfg
}
