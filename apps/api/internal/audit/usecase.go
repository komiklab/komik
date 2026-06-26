package audit

import (
	"github.com/komiklab/komik/internal/client"
	"github.com/komiklab/komik/internal/models"
	"github.com/komiklab/komik/internal/repositories"
)

// AuditService exposes use-case operations for the audit-log domain.
type AuditService struct {
	repo *repositories.AuditLogRepository
}

// NewAuditService creates a new AuditService backed by the given Postgres client.
func NewAuditService(dbcon *client.PostgresClient) *AuditService {
	return &AuditService{
		repo: repositories.NewAuditLogRepository(dbcon),
	}
}

// GetAuditLogs returns a paginated list of audit log entries together with the
// total row count, delegating directly to the repository layer.
func (a *AuditService) GetAuditLogs(limit, offset int) ([]models.AuditLog, int64, error) {
	return a.repo.GetAuditLogs(limit, offset)
}

func (a *AuditService) SaveAuditLog(auditLog *models.AuditLog) error {
	return a.repo.Save(auditLog)
}