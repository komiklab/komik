package main

import (
	//"os"

	"os"
	"os/signal"
	"syscall"

	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/controller"
	"github.com/komiklab/komik/internal/event"
	httpserver "github.com/komiklab/komik/internal/httpServer"
	"github.com/komiklab/komik/internal/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := internal.NewConfig()
	if cfg.IsDebugLoggerConfig {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.With().Caller().Logger()
	postgresclient := client.NewPostgresClient(cfg)
	postgresclient.Migrate(&models.Admin{})
	contrlr := controller.NewController(cfg)
	publisher := event.NewPublisher(cfg)
	httpComponent := httpserver.NewHttpComponent(cfg, postgresclient, publisher)
	eventComponent := event.NewEventLoopComponent(cfg, postgresclient)
	contrlr.AddComponent(eventComponent)
	contrlr.AddComponent(httpComponent)
	contrlr.Init()
	contrlr.Start()
	stopchannel := make(chan os.Signal, 1)
	signal.Notify(stopchannel, os.Interrupt, syscall.SIGTERM)
	receivedSignal := <-stopchannel
	log.Info().Msg("Received signal " + receivedSignal.String())
	contrlr.Stop()
}
