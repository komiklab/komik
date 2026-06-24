package repositories

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/utils"
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

// GetAuditLogs returns a paginated slice of audit log entries plus the total
// count of all rows in the table (for pagination metadata).
func (a *AuditLogRepository) GetAuditLogs(limit, offset int) ([]models.AuditLog, int64, error) {
	gormdb := a.dbcon.GetClient()

	var total int64
	if err := gormdb.Model(&models.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, utils.NewDatabaseError("failed to count audit logs", err)
	}

	var logs []models.AuditLog
	if err := gormdb.
		Order("occurred_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, utils.NewDatabaseError("failed to fetch audit logs", err)
	}

	return logs, total, nil
}
