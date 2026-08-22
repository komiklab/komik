package event

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal/audit"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/entity"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/orchestrator"


	//"github.com/komiklab/komik/internal/repositories"

	"github.com/rs/zerolog/log"
)

type EventHandler struct {
	//auditLogRepo *repositories.AuditLogRepository
	auditSrv           *audit.AuditService
	entitySrv          *entity.EntityService
	orchestratorClient *orchestrator.OrchestratorClient
}

func NewEventHandler(dbcon *client.PostgresClient, tempoClient *orchestrator.OrchestratorClient) *EventHandler {
	//auditLogRepo := repositories.NewAuditLogRepository(dbcon)
	authsrv := audit.NewAuditService(dbcon)
	entitySrv := entity.NewEntityService(dbcon)
	return &EventHandler{
		auditSrv:           authsrv,
		entitySrv:          entitySrv,
		orchestratorClient: tempoClient,
	}
}

func (h *EventHandler) HandleAuditLog(msg *message.Message) error {
	log.Info().Msg("Received Audit log handle.")
	auditlogData := msg.Payload
	auditlogModel := &models.AuditLog{}
	err := auditlogModel.Unmarshal(auditlogData)
	// if utils.IsErrNotNil(err) {
	// 	return err
	// }
	//err = h.auditLogRepo.Save(auditlogModel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal audit log")
		return err
	}
	err = h.auditSrv.SaveAuditLog(auditlogModel)
	// if utils.IsErrNotNil(err) {
	// 	return err
	// }
	if err != nil {
		log.Error().Err(err).Msg("Failed to save audit log")
		return err
	}
	return nil
}

func (h *EventHandler) HandleEntityInitiated(msg *message.Message) error {
	// log.Info().Msg("Received Entity log handle.")
	// log.Info().Msg("causation id should be " + msg.UUID)
	ctx := context.Background()
	auditlogData := msg.Payload
	auditlogModel := &models.AuditLog{}
	err := auditlogModel.Unmarshal(auditlogData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal audit log")
		return err
	}
	entityId := auditlogModel.Data
	entityData, err := h.entitySrv.RetrieveEntityById(msg.Context(), entityId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve entity")
		return err
	}
	transitioner := entity.NewEntityTransitioner(entityData)
	err = transitioner.Dispatch(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to change the state")
		return err
	}
	resp, err := h.orchestratorClient.SendEvent(ctx, orchestrator.INNGEST_ENTITY_DISPATCH_EVENT, entityId, entityData.StructToMap())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send event to orchestrator")
		return err
	}
	log.Debug().Msg("Event sent to orchestrator " + resp)
	entityData.AgentThreadId = &resp
	entityData.Status = entity.EntityStateDispatched
	err = h.entitySrv.Update(entityData, map[string]string{"causationId": msg.UUID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update entity")
		return err
	}
	return nil
}
