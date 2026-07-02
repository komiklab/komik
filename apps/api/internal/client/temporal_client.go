package client

import (
	"context"
	"time"

	"github.com/komiklab/komik/internal"
	"github.com/rs/zerolog/log"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type TemporalClient struct {
	Client client.Client
	TaskQueue string
}

func (t *TemporalClient) GetClient() client.Client {
	return t.Client
}

func (t *TemporalClient) StartWorkflow(ctx context.Context, workflowID string, workflowName string, workflowArgs ...any)(client.WorkflowRun, error) {
	workflowOpt := client.StartWorkflowOptions{
		ID: workflowID,
		TaskQueue: t.TaskQueue,
		WorkflowExecutionTimeout: 24*time.Hour,
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}
	return t.Client.ExecuteWorkflow(ctx, workflowOpt, workflowName, workflowArgs...)
}

func (t *TemporalClient) Close() {
	t.Client.Close()
}

func NewTemporalClient(cfg *internal.Config) *TemporalClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Tempclient, err := client.DialContext(ctx, client.Options{
		Namespace: cfg.TemporalConfig.Namespace,
		HostPort:  cfg.TemporalConfig.TemporalUrl,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to temporal")
	}
	return &TemporalClient{
		Client: Tempclient,
		TaskQueue: cfg.TemporalConfig.TaskQueue,
	}
}