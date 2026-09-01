package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"github.com/google/uuid"
)

type Entity struct {
	Id                 uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt          time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	CompletedAt        *time.Time      `json:"completed_at" gorm:"default:null"`
	SourceType         string          `json:"source_type" gorm:"type:varchar(255);not null"`
	SourceRef          string          `json:"source_ref" gorm:"type:varchar(255);not null;uniqueIndex"`
	SourcePayload      datatypes.JSON  `json:"source_payload" gorm:"type:jsonb;not null;"`
	Status             string          `json:"status" gorm:"type:varchar(20);not null;default:'initiated'"`
	AgentThreadId      *string         `json:"agent_thread_id" gorm:"type:varchar(255);default:null"`
	ActiveInterruptId  *uuid.UUID      `json:"active_interrupt_id" gorm:"type:uuid;default:null"`
	ResultSummary      *string         `json:"result_summary" gorm:"type:text;default:null"`
	Result             *datatypes.JSON `json:"result" gorm:"type:jsonb;default:null"`
	ErrorMessage       *string         `json:"error_message" gorm:"type:text;default:null"`
}

func (e *Entity) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Entity) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *Entity) StructToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                 e.Id,
		"created_at":         e.CreatedAt,
		"updated_at":         e.UpdatedAt,
		"completed_at":       e.CompletedAt,
		"source_type":        e.SourceType,
		"source_ref":         e.SourceRef,
		"source_payload":     e.SourcePayload,
		"status":             e.Status,
		"agent_thread_id":      e.AgentThreadId,
		"active_interrupt_id":  e.ActiveInterruptId,
		"result_summary":     e.ResultSummary,
		"result":             e.Result,
		"error_message":      e.ErrorMessage,
	}
}

