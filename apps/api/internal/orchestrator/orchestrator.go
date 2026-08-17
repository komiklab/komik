package orchestrator

import (
	"context"

	"github.com/inngest/inngestgo"
	"github.com/inngest/inngestgo/connect"
	"github.com/komiklab/komik/internal"
	"github.com/rs/zerolog/log"
)

type Orchestrator struct {
	cfg    *internal.Config
	client inngestgo.Client
	conn   connect.WorkerConnection
}

func NewOrchestrator(cfg *internal.Config) *Orchestrator {
	println(cfg.InngestConfig.EventKey)
	client, err := inngestgo.NewClient(inngestgo.ClientOpts{
		AppID:      cfg.InngestConfig.AppID,
		EventKey:   inngestgo.StrPtr(cfg.InngestConfig.EventKey),
		SigningKey: inngestgo.StrPtr(cfg.InngestConfig.SigningKey),
		AppVersion: inngestgo.StrPtr(cfg.InngestConfig.AppVersion),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Inngest client")
	}
	return &Orchestrator{
		cfg:    cfg,
		client: client,
	}
}

func (o *Orchestrator) Start() {
	conn, err := inngestgo.Connect(context.Background(), inngestgo.ConnectOpts{
		Apps: []inngestgo.Client{o.client},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Inngest")
	}
	o.conn = conn
}


func (o *Orchestrator) Stop() error {
	return o.conn.Close()
}
