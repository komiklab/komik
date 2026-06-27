package event

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal/audit"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"

	//"github.com/komiklab/komik/internal/repositories"

	"github.com/rs/zerolog/log"
)

type EventHandler struct {
	//auditLogRepo *repositories.AuditLogRepository
	auditSrv *audit.AuditService
}

func NewEventHandler(dbcon *client.PostgresClient) *EventHandler {
	//auditLogRepo := repositories.NewAuditLogRepository(dbcon)
	authsrv := audit.NewAuditService(dbcon)
	return &EventHandler{
		auditSrv: authsrv,
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

func (h *EventHandler) HandleEntity(msg *message.Message) error {
	log.Info().Msg("Received Entity log handle.")
	log.Info().Msg("causation id should be " + msg.UUID)
	entityData := msg.Payload
	log.Info().Msg("entity data is " + string(entityData))
	return nil
}
