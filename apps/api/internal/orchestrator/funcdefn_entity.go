package orchestrator

import (
	"context"

	"github.com/inngest/inngestgo"
	"github.com/komiklab/komik/internal/models"
	"github.com/rs/zerolog/log"
)

func (o *Orchestrator) registerEntityTransitionFn() {
	inngestgo.CreateFunction(o.client,
		inngestgo.FunctionOpts{
			ID:   "handling-entity-transition",
			Name: "Handling Entity Transition",
		},
		inngestgo.EventTrigger("entity/dispatched", nil),
		func(ctx context.Context, input inngestgo.Input[models.Entity]) (any, error) {
			log.Debug().Msgf("received entity %v", input)
			return nil, nil
		}, 
	)
}
