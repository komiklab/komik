package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
	taskqueue "github.com/komiklab/komik/internal/task_queue"
	"github.com/rs/zerolog/log"
)

type WorkerComponent struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	cfg *internal.Config
}

func NewWorkerComponent(redisconn *client.RedisClient, cfg *internal.Config) *WorkerComponent {
	server := asynq.NewServerFromRedisClient(redisconn.GetClient(), asynq.Config{})
	return &WorkerComponent{
		cfg: cfg,
		server: server,
	}
}

// GetName implements [internal.Component].
func (w *WorkerComponent) GetName() string {
	return "AsyncWorkerComponent"
}

// Init implements [internal.Component].
func (w *WorkerComponent) Init() {
	handler := NewWorkerHandler(w.cfg)
	mux := asynq.NewServeMux()
	mux.HandleFunc(string(taskqueue.SlackAppMention), handler.ProcessTask)
	w.mux = mux
}

// Start implements [internal.Component].
func (w *WorkerComponent) Start() {
	err := w.server.Run(w.mux)
	if err != nil {
		log.Fatal().Err(err).Msg("Worker failed to start")
	}
}

// Stop implements [internal.Component].
func (w *WorkerComponent) Stop(ctx context.Context) {
	w.server.Stop()
	w.server.Shutdown()
}

var _ internal.Component = (*WorkerComponent)(nil)
