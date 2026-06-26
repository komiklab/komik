package event

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal/audit"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"

	//"github.com/komiklab/komik/internal/repositories"
	"github.com/komiklab/komik/internal/utils"
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
	if utils.IsErrNotNil(err) {
		return err
	}
	//err = h.auditLogRepo.Save(auditlogModel)
	err = h.auditSrv.SaveAuditLog(auditlogModel)
	if utils.IsErrNotNil(err) {
		return err
	}
	return nil
}

func (h *EventHandler) HandleEntity(msg *message.Message) error {
	log.Info().Msg("Received Entity log handle.")
	entityData := msg.Payload
	entityModel := &models.Entity{}
	err := entityModel.Unmarshal(entityData)
	if err != nil {
		return err
	}
	return errors.New("for fun in HandleEntity")
}
