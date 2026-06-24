package httphandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/audit"
	"github.com/komiklab/komik/internal/utils"
	"github.com/labstack/echo/v5"
)

// GetAuditlog implements [apidefn.ServerInterface].
func (h *HttpHandler) GetAuditlog(ctx *echo.Context, params apidefn.GetAuditlogParams) error {
	svc := audit.NewAuditService(h.dbclient)

	logs, total, err := svc.GetAuditLogs(params.Limit, params.Offset)
	if err != nil {
		komikErr, ok := err.(*utils.KomikError)
		if !ok {
			komikErr = utils.NewGeneralError(err)
		}
		return ctx.JSON(int(komikErr.StatusCode), komikErr)
	}

	items := make([]apidefn.AuditlogResponse, 0, len(logs))
	for _, l := range logs {
		l := l // capture loop var

		eventId := uuid.UUID(l.EventId)
		occurredAt := time.Unix(l.OccurredAt, 0).UTC()

		items = append(items, apidefn.AuditlogResponse{
			EventId:       &eventId,
			EventVersion:  &l.EventVersion,
			OccurredAt:    &occurredAt,
			CorrelationId: &l.CorrelationId,
			InitiatorId:   &l.InitiatorId,
			InitiatorType: &l.InitiatorType,
			ResourceType:  &l.ResourceType,
			Severity:      &l.Severity,
			Data:          &l.Data,
			EventType:     &l.EventType,
		})
	}

	totalInt := int(total)
	limit := params.Limit
	offset := params.Offset

	resp := apidefn.AuditlogGetResponse{
		Items: &items,
		Metadata: &struct {
			Limit  *int `json:"limit,omitempty"`
			Offset *int `json:"offset,omitempty"`
			Total  *int `json:"total,omitempty"`
		}{
			Limit:  &limit,
			Offset: &offset,
			Total:  &totalInt,
		},
	}

	_ = fmt.Sprintf // ensure fmt is used via logging if needed
	return ctx.JSON(http.StatusOK, resp)
}
