package orchestrator

import (
	"context"

	"github.com/inngest/inngestgo"
	"github.com/rs/zerolog/log"
)

func (o *Orchestrator) RegisterFuncs() {
	log.Debug().Msg("registering functions")
	// register all the functions
	inngestgo.CreateFunction(o.client,
		inngestgo.FunctionOpts{
			ID:   "handling-entity-transition",
			Name: "Handling Entity Transition",
		},
		inngestgo.EventTrigger("entity/dispatched", nil),
		func(ctx context.Context, input inngestgo.Input[map[string]any]) (any, error) {
			println("received")
			return nil, nil
		},
	)

}
