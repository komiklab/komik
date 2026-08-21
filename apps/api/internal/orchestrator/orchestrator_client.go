package orchestrator

import (
	"context"

	"github.com/inngest/inngestgo"
)

type OrchestratorClient struct {
	client inngestgo.Client
}

func NewOrchestratorClient(o *Orchestrator) *OrchestratorClient {
	return &OrchestratorClient{
		client: o.client,
	}
}

func (oc *OrchestratorClient)  SendEvent(ctx context.Context, name string, id string, data map[string]any) (string, error) {
	resp, err := oc.client.Send(ctx, inngestgo.Event{
		Name: name,
		ID:   inngestgo.StrPtr(id),
		Data: data,
	})
	return resp, err
}
