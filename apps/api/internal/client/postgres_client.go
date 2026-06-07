package client

import (
	"github.com/komiklab/komik/internal"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ Client = (*PostgresClient)(nil)

type PostgresClient struct {
	gormdb *gorm.DB
}

func (p *PostgresClient) GetClient() any {
	return p.gormdb
}

func (p *PostgresClient) Ping() error {
	db, err := p.gormdb.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database connection")
		return err
	}
	return db.Ping()
}

func (p *PostgresClient) Close() {
	db, err := p.gormdb.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database connection")
		return
	}
	db.Close()
}

func (p *PostgresClient) Migrate(models ...any) {
	err := p.gormdb.AutoMigrate(models...)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Debug().Msg("Database migrated")
}

func NewPostgresClient(cfg *internal.Config) *PostgresClient {
	dsn := cfg.PostgresDSN
	gormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	db, err := gormdb.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database connection")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}
	log.Debug().Msg("Connected to database")
	return &PostgresClient{
		gormdb: gormdb,
	}
}
