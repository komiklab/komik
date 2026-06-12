package event

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

type EventHandler struct {
	auditLogRepo *repositories.AuditLogRepository
}

func NewEventHandler(dbcon *client.PostgresClient) *EventHandler {
	auditLogRepo := repositories.NewAuditLogRepository(dbcon)
	return &EventHandler{
		auditLogRepo: auditLogRepo,
	}
}

func (h *EventHandler) HandleAuditLog(msg *message.Message) error {
	auditlogData := msg.Payload
	var auditlogModel *models.AuditLog
	err := auditlogModel.Unmarshal(auditlogData)
	if err != nil {
		return err
	}
	return h.auditLogRepo.Save(auditlogModel)
}
