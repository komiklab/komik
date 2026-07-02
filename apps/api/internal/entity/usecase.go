package entity

import (
	"context"
	"errors"
	"time"

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

func (e *EntityService) Update(entityData *models.Entity, metadata map[string]string) error {
	msg, err := e.createMsgEnvelope(entityData, utils.SystemInitiator, metadata, EventTypeEntityDispatched)
	if err != nil {
		log.Error().Err(err).Msg("failed to create msg envelope")
		return err
	}
	return e.repo.Save(entityData, msg, utils.DispatchedTopic)
}

func (e *EntityService) RetrieveEntityById(ctx context.Context, id string) (*models.Entity, error) {
	return e.repo.GetEntityById(ctx, id)
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
	envelope, err := e.createMsgEnvelope(entity, initiator, nil, EventTypeEntityInitiated)
	if err != nil {
		log.Error().Err(err).Msg("failed to create msg envelope")
		return err
	}
	return e.repo.Save(entity, envelope, utils.InitiatedTopic)
}

func NewEntityService(dbcon *client.PostgresClient) *EntityService {
	return &EntityService{
		repo: repositories.NewEntityRepo(dbcon),
	}
}

func (e *EntityService) createMsgEnvelope(entity *models.Entity, initiator string, metadata map[string]string, eventtype string) (*message.Message, error) {
	eventID := uuid.New()
	causationID := entity.SourceRef
	// if metadata is not nil and contains a field called causation id we should use that instead of entity.SourceRef
	if metadata != nil {
		if caus, ok := metadata["causationId"]; ok {
			causationID = caus
		}
	}
	envelope := &models.AuditLog{
		EventId:       eventID,
		EventType:     eventtype,
		OccurredAt:    time.Now().UnixMilli(),
		EventVersion:  1,
		CorrelationId: entity.SourceRef,
		InitiatorId:   initiator,
		InitiatorType: entity.SourceType,
		ResourceType:  entity.Id.String(),
		Severity:      "INFO",
		Payload:       string(entity.SourcePayload),
		Data:          entity.Id.String(),
		CausationId:   causationID,
	}
	envelopeBytes, err := envelope.Marshal()
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal envelope")
		return nil, err
	}
	msg := message.NewMessage(eventID.String(), envelopeBytes)
	middleware.SetCorrelationID(envelope.CorrelationId, msg)
	return msg, nil
}
