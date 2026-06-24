package auth

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/google/uuid"
	"github.com/komiklab/komik/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	EventAdminCreatedSubject         = "komik.admin.created"
	EventAdminCreationFaailedSubject = "komik.admin.creation_failed"
)

const (
	EventAdminCreated        = "AdminCreated"
	EventAdminCreationFailed = "AdminCreationFailed"
)

const (
	InitiatorTypeAdmin = "Admin"
)

func (a *AuthService) MessageEventAdminCreated(admin string, correlationId string) (*message.Message, error) {
	eventId, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create event ID")
		return nil, err
	}
	eventData := models.AuditLog{
		EventId:       eventId,
		EventType:     EventAdminCreated,
		EventVersion:  1,
		CorrelationId: correlationId,
		InitiatorId:   admin,
		InitiatorType: InitiatorTypeAdmin,
		ResourceType:  InitiatorTypeAdmin,
		Severity:      "Info",
	}
	jsonData, err := eventData.Marshal()
	if err != nil {
		return nil, err
		// ignore the error
	}
	msg := message.NewMessage(watermill.NewUUID(), jsonData)
	middleware.SetCorrelationID(watermill.NewUUID(), msg)
	return msg, nil
}
