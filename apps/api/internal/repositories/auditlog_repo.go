package repositories

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"gorm.io/gorm/clause"
)

type AuditLogRepository struct {
	dbcon *client.PostgresClient
}

func NewAuditLogRepository(dbcon *client.PostgresClient) *AuditLogRepository {
	return &AuditLogRepository{
		dbcon: dbcon,
	}
}

func (a *AuditLogRepository) Save(auditLog *models.AuditLog) error {
	gormdb := a.dbcon.GetClient()
	return gormdb.Clauses(clause.OnConflict{DoNothing: true}).Create(auditLog).Error
}
