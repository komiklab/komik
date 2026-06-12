package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AuditLog struct {
	EventId       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventType     string    `gorm:"type:varchar(255);not null"`
	EventVersion  int       `gorm:"not null;default:1"`
	OccurredAt    int64     `gorm:"autoCreateTime"`
	CorrelationId string    `gorm:"type:string;not null"`
	InitiatorId   string    `gorm:"type:varchar(255);not null"`
	InitiatorType string    `gorm:"type:varchar(255);not null"`
	ResourceType  string    `gorm:"type:varchar(255);not null"`
	Severity      string    `gorm:"type:varchar(255);not null"`
	Payload       string    `gorm:"serializer:json;not null"`
	Data          string    `gorm:"serializer:json;not null"`
}

func (a *AuditLog) Unmarshal(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *AuditLog) Marshal() ([]byte, error) {
	return json.Marshal(a)
}
