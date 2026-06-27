package entity

import (
	"context"
	"errors"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/google/uuid"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
	"github.com/komiklab/komik/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
)

type EntityService struct {
	repo *repositories.EntityRepo
}

func (e *EntityService) InitiateEntity(ctx context.Context, sourceType, sourceRef, initiator string, sourcePayload []byte) error {
	// first validate the sourcerType
	if !utils.IsValidSourceType(sourceType) {
		return errors.New("Not a valid sourceType: " + sourceType)
	}
	entity := &models.Entity{
		Id:            datatypes.NewUUIDv4(),
		Status:        EntityStateInitiated,
		SourceType:    sourceType,
		SourceRef:     sourceRef,
		SourcePayload: sourcePayload,
	}
	envelope, err := e.CreateMsgEnvelope(entity, initiator)
	if err != nil {
		log.Error().Err(err).Msg("failed to create msg envelope")
		return err
	}
	return e.repo.Save(entity, envelope)
}

func NewEntityService(dbcon *client.PostgresClient) *EntityService {
	return &EntityService{
		repo: repositories.NewEntityRepo(dbcon),
	}
}

func (e *EntityService) CreateMsgEnvelope(entity *models.Entity, initiator string) (*message.Message, error) {
	// entityBytes, err := entity.Marshal()
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to marshal entity")
	// 	return nil, err
	// }
	envelope := &models.AuditLog{
		EventId:       uuid.New(),
		EventType:     EventTypeEntityInitiated,
		OccurredAt:    time.Now().UnixMilli(),
		EventVersion:  1,
		CorrelationId: entity.SourceRef,
		InitiatorId:   initiator,
		InitiatorType: entity.SourceType,
		ResourceType:  entity.Id.String(),
		Severity:      "INFO",
		Payload:       string(entity.SourcePayload),
		Data:          entity.Id.String(),
		CausationId:   entity.SourceRef,
	}
	envelopeBytes, err := envelope.Marshal()
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal envelope")
		return nil, err
	}
	msg := message.NewMessage(watermill.NewUUID(), envelopeBytes)
	middleware.SetCorrelationID(envelope.CorrelationId, msg)
	return msg, nil
}
