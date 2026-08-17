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
	"github.com/komiklab/komik/internal/orchestrator"
	"github.com/komiklab/komik/internal/task_queue/worker"
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
	// we will set colorful os.std
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Debug().Msg("Starting API server with debug logging enabled ")
	postgresclient := client.NewPostgresClient(cfg)
	postgresclient.Migrate(
		models.Admin{},
		models.AuditLog{},
		models.UserRepresentation{},
		models.Entity{},
		models.EntityInterrupt{},
		models.Agent{},
		models.Hooks{},
	)
	redisClient := client.NewRedisClient(cfg)
	//temporalClient := client.NewTemporalClient(cfg)
	contrlr := controller.NewController(cfg)
	publisher := event.NewPublisher(cfg)
	orchestratorComponent := orchestrator.NewOrchestratorComponent(cfg)
	orchestratorClient := orchestratorComponent.GetOrchestratorClient()
	httpComponent := httpserver.NewHttpComponent(cfg, postgresclient, publisher, redisClient)
	eventComponent := event.NewEventLoopComponent(cfg, postgresclient, orchestratorClient)
	workerComponent := worker.NewWorkerComponent(redisClient, postgresclient, cfg)
	
	contrlr.AddComponent(eventComponent)
	contrlr.AddComponent(httpComponent)
	contrlr.AddComponent(workerComponent)
	contrlr.AddComponent(orchestratorComponent)
	contrlr.Init()
	contrlr.Start()
	stopchannel := make(chan os.Signal, 1)
	signal.Notify(stopchannel, os.Interrupt, syscall.SIGTERM)
	receivedSignal := <-stopchannel
	log.Info().Msg("Received signal " + receivedSignal.String())
	contrlr.Stop()
}
