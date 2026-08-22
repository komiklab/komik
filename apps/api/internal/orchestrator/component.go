package orchestrator

import (
	"context"

	"github.com/komiklab/komik/internal"
	"github.com/komiklab/komik/internal/client"
)

type OrchestratorComponent struct {
	orch *Orchestrator
}

func NewOrchestratorComponent(cfg *internal.Config, dbcon *client.PostgresClient) *OrchestratorComponent {
	return &OrchestratorComponent{
		orch: NewOrchestrator(cfg, dbcon),
	}
}

// GetName implements Component.
func (w *OrchestratorComponent) GetName() string {
	return "OrchestratorComponent"
}

// Init implements Component.
func (w *OrchestratorComponent) Init() {
	w.orch.RegisterFuncs()
}

// Start implements Component.
func (w *OrchestratorComponent) Start() {
	w.orch.Start()
}

// Stop implements Component.
func (w *OrchestratorComponent) Stop(ctx context.Context) {
	w.orch.Stop()
}

func (w *OrchestratorComponent) GetOrchestrator() *Orchestrator {
	return w.orch
}

func (w *OrchestratorComponent) GetOrchestratorClient() *OrchestratorClient {
	return NewOrchestratorClient(w.orch)
}

var _ internal.Component = (*OrchestratorComponent)(nil)
