package orchestrator

import (
	"context"

	"github.com/inngest/inngestgo"
	"github.com/inngest/inngestgo/step"
	"github.com/komiklab/komik/internal/conversion"
	"github.com/komiklab/komik/internal/models"
	"github.com/rs/zerolog/log"
)

func (o *Orchestrator) registerEntityTransitionFn() {
	inngestgo.CreateFunction(o.client,
		inngestgo.FunctionOpts{
			ID:   "handling-entity-transition",
			Name: "Handling Entity Transition",
		},
		inngestgo.EventTrigger(INNGEST_ENTITY_DISPATCH_EVENT, nil),
		func(ctx context.Context, input inngestgo.Input[models.Entity]) (any, error) {
			log.Debug().Msgf("received entity %v", input)
			// Step 1. Fetche the conversation based on session ID
			entity := input.Event.Data
			sessionId := entity.Id
			convSrv := conversion.NewConverstionSrv(o.dbcon)
			conversation, err := step.Run(ctx, "fetchConversationBySessionId", func(ctx context.Context) ([]*models.Conversation, error) {
				return convSrv.GetConversationBySessionId(sessionId)
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch conversation")
				return nil, err
			}
			if len(conversation) == 0 {
				// when there is no conversation, we should create one
				// note this means it comes from hook
				err := convSrv.CreateConversation(&models.Conversation{
					SessionId:        sessionId,
					OwnerId:          input.InputCtx.RunID,
					ConversationType: conversion.ConversationTypeHook,
					Sequence:         1,
					Content:          entity.SourcePayload,
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed to create conversation")
					return nil, err
				}
			}
			// step 2: call needle to get the list of tools
			return nil, nil
		},
	)
}
