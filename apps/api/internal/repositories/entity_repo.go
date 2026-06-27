package repositories

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EntityRepo struct {
	dbcon  *client.PostgresClient
	fwdpub *ForwarderPublisher
}

func (e *EntityRepo) Save(entity *models.Entity, envelope *message.Message) error {
	gormdb := e.dbcon.GetClient()
	return gormdb.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "source_ref"},
			},
			DoNothing: true,
		}).Create(entity)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			log.Warn().Msgf("Entity with Id %s, source ref %s already exists", entity.Id.String(), entity.SourceRef)
			return nil
		}
		return e.fwdpub.Publish(tx, "komik.entity.initiated", envelope)
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
