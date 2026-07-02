package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AuditLog struct {
	EventId       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"event_id"`
	EventType     string    `gorm:"type:varchar(255);not null" json:"event_type"`
	EventVersion  int       `gorm:"not null;default:1" json:"event_version"`
	OccurredAt    int64     `gorm:"autoCreateTime" json:"occurred_at"`
	CorrelationId string    `gorm:"type:string;not null" json:"correlation_id"`
	InitiatorId   string    `gorm:"type:varchar(255);not null" json:"initiator_id"`
	InitiatorType string    `gorm:"type:varchar(255);not null" json:"initiator_type"`
	ResourceType  string    `gorm:"type:varchar(255);not null" json:"resource_type"`
	Severity      string    `gorm:"type:varchar(255);not null" json:"severity"`
	Payload       string    `gorm:"serializer:json;not null" json:"payload"`
	Data          string    `gorm:"serializer:json;not null" json:"data"`
	CausationId   string    `gorm:"type:string;default:'';not null" json:"causation_id"`
}

func (a *AuditLog) Unmarshal(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *AuditLog) Marshal() ([]byte, error) {
	return json.Marshal(a)
}
