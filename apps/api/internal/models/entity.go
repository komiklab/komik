package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type Entity struct {
	Id                 datatypes.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt          time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	CompletedAt        *time.Time      `json:"completed_at" gorm:"default:null"`
	SourceType         string          `json:"source_type" gorm:"type:varchar(255);not null"`
	SourceRef          string          `json:"source_ref" gorm:"type:varchar(255);not null;uniqueIndex"`
	SourcePayload      datatypes.JSON  `json:"source_payload" gorm:"type:jsonb;not null;"`
	Status             string          `json:"status" gorm:"type:varchar(20);not null;default:'initiated'"`
	TemporalWorkflowId *string         `json:"temporal_workflow_id" gorm:"type:varchar(255);default:null"`
	TemporalRunId      *string         `json:"temporal_run_id" gorm:"type:varchar(255);default:null"`
	TemporalTaskQueue  string          `json:"temporal_task_queue" gorm:"type:varchar(255);default:'default'"`
	AgentThreadId      *string         `json:"agent_thread_id" gorm:"type:varchar(255);default:null"`
	ActiveInterruptId  *datatypes.UUID `json:"active_interrupt_id" gorm:"type:uuid;default:null"`
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
		"temporal_workflow_id": e.TemporalWorkflowId,
		"temporal_run_id":      e.TemporalRunId,
		"temporal_task_queue":  e.TemporalTaskQueue,
		"agent_thread_id":      e.AgentThreadId,
		"active_interrupt_id":  e.ActiveInterruptId,
		"result_summary":     e.ResultSummary,
		"result":             e.Result,
		"error_message":      e.ErrorMessage,
	}
}

type EntityInterrupt struct {
	Id              datatypes.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	EntityID        datatypes.UUID `json:"entity_id" gorm:"type:uuid;not null;index"`
	Entity          Entity         `json:"entity" gorm:"foreignKey:EntityID"`
	InterruptType   string         `json:"interrupt_type" gorm:"type:varchar(20);not null"`
	InterruptReason string         `json:"interrupt_reason" gorm:"type:varchar(255);not null"`
	Reason          string         `json:"reason" gorm:"type:text;not null"`
	ContextSnapShot datatypes.JSON `json:"context_snapshot" gorm:"type:jsonb;not null"`
	PromptSend      string         `json:"prompt_send" gorm:"type:text;not null"`
	InterruptStatus string         `json:"interrupt_status" gorm:"type:varchar(20);not null"`
	ResponsePayload datatypes.JSON `json:"response_payload" gorm:"type:jsonb;not null"`
	ResolvedBy      string         `json:"resolved_by" gorm:"type:text;not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	RespondBy       time.Time      `json:"respond_by" gorm:"not null"`
	ResolvedAt      time.Time      `json:"resolved_at" gorm:"default:null"`
}
