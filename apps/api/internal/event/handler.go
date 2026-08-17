package event

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal/audit"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/entity"
	"github.com/komiklab/komik/internal/models"

	//"github.com/komiklab/komik/internal/repositories"

	"github.com/rs/zerolog/log"
)

type EventHandler struct {
	//auditLogRepo *repositories.AuditLogRepository
	auditSrv       *audit.AuditService
	entitySrv      *entity.EntityService
	temporalClient *client.TemporalClient
}

func NewEventHandler(dbcon *client.PostgresClient, tempoClient *client.TemporalClient) *EventHandler {
	//auditLogRepo := repositories.NewAuditLogRepository(dbcon)
	authsrv := audit.NewAuditService(dbcon)
	entitySrv := entity.NewEntityService(dbcon)
	return &EventHandler{
		auditSrv:       authsrv,
		entitySrv:      entitySrv,
		temporalClient: tempoClient,
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
	// temporalDetails, err := h.temporalClient.StartWorkflow(ctx, entityId, "EntityTransitionWorkflow", entityData)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to start temporal workflow")
	// 	return err
	// }
	// log.Info().Msg("Temporal workflow started " + temporalDetails.GetRunID())
	// workflowId := temporalDetails.GetID()
	// workflowRunID := temporalDetails.GetRunID()
	// entityData.TemporalRunId = &workflowRunID
	// entityData.TemporalWorkflowId = &workflowId
	// entityData.TemporalTaskQueue = h.temporalClient.TaskQueue
	// entityData.Status = entity.EntityStateDispatched
	// err = h.entitySrv.Update(entityData, map[string]string{"causationId": msg.UUID})
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to update entity")
	// 	return err
	// }
	return nil
}
