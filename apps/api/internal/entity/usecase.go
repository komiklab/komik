package entity

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

type EntityService struct {
	repo *repositories.EntityRepo
}

func (e *EntityService) InitiateEntity(sourceType, sourceRef, initiator string, sourcePayload []byte) error {
	entity := &models.Entity{
		Status:        EntityStateInitiated,
		SourceType:    sourceType,
		SourceRef:     sourceRef,
		SourcePayload: sourcePayload,
	}
	return e.repo.Save(entity, initiator)
}

func NewEntityService(dbcon *client.PostgresClient) *EntityService {
	return &EntityService{
		repo: repositories.NewEntityRepo(dbcon),
	}
}


