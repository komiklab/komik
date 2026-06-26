package repositories

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"gorm.io/gorm"
)

type EntityRepo struct {
	dbcon  *client.PostgresClient
	fwdpub *ForwarderPublisher
}

func (e *EntityRepo) Save(entity *models.Entity, inititator string) error {
	gormdb := e.dbcon.GetClient()
	return gormdb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		databytes, err := entity.Marshal()
		if err != nil {
			return err
		}
		eventData := models.AuditLog{
			EventType:     "EntityInitiated",
			CorrelationId: entity.SourceRef,
			InitiatorId:   inititator,
			InitiatorType: "user",
			ResourceType:  entity.SourceType,
			Severity:      "Info",
			Payload:       string(entity.SourcePayload),
			Data:          string(databytes),
		}
		jsonData, err := eventData.Marshal()
		if err != nil {
			return err
		}
		msg := message.NewMessage(watermill.NewUUID(), jsonData)
		middleware.SetCorrelationID(watermill.NewUUID(), msg)
		return e.fwdpub.Publish(tx, "komik.ent.initiated", msg)
	})
}

func (e *EntityRepo) GetByID(id string) (*models.Entity, error) {
	gormdb := e.dbcon.GetClient()
	var entity models.Entity
	err := gormdb.Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (e *EntityRepo) Update(entity *models.Entity) error {
	gormdb := e.dbcon.GetClient()
	return gormdb.Save(entity).Error
}

func NewEntityRepo(dbcon *client.PostgresClient) *EntityRepo {
	return &EntityRepo{
		dbcon:  dbcon,
		fwdpub: NewForwarderPublisher(),
	}
}
