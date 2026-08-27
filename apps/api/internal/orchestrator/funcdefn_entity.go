package orchestrator

import (
	"context"

	"github.com/inngest/inngestgo"
	"github.com/inngest/inngestgo/step"
	"github.com/komiklab/komik/internal/agent"
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
			agentsrv := agent.NewAgentService(o.dbcon)
			conversation, err := step.Run(ctx, "fetchConversationBySessionId", func(ctx context.Context) ([]*models.Conversation, error) {
				conversation, err := convSrv.GetConversationBySessionId(sessionId)
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
				return conversation, nil
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch conversation")
				return nil, err
			}
			log.Debug().Msgf("fetched conversation %v", conversation)

			// step 2: get the agents based on the entity source type
			agents, err := step.Run(ctx, "fetchAgentsBasedOnEntity", func(ctx context.Context) ([]models.Agent, error) {
				entity := input.Event.Data
				agentsList, err := agentsrv.ListAgentBasedOnEntity(entity)
				if err != nil {
					log.Error().Err(err).Msg("Failed to fetch agents")
					return nil, err
				}
				agents := agentsList.Agents
				log.Debug().Msgf("fetched agents %v", agents)
				return agents, nil
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch agents")
				return nil, err
			}
			log.Debug().Msgf("fetched agents %v", agents)
			// TODO: logic to decide when there are no agents and multiple agent
			// TODO: as of now we will consider one agent only
			// step 3: call the agents
			agent := agents[0]
			_, err = step.Run(ctx, "callAgent", func(ctx context.Context) (any, error) {
				// call the agent
				_, err := agentsrv.CallAgent(agent, entity)
				if err != nil {
					log.Error().Err(err).Msg("Failed to call agent")
					return nil, err
				}
				return nil, nil
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to call agent")
				return nil, err
			}
			log.Debug().Msgf("called agent %v", agent)
			return nil, nil

		},
	)
}
