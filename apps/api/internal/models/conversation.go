package models

import (
	"time"

	"gorm.io/datatypes"
	"github.com/google/uuid"
)

type Conversation struct {
	Id uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	SessionId        uuid.UUID      `json:"session_id" gorm:"type:uuid;required"`
	OwnerId          string         `json:"user_id" gorm:"type:varchar(255);not null;default:required"`
	ConversationType string         `json:"conversation_type" gorm:"type:varchar(20);not null;default:'user'"`
	Sequence         int64          `json:"sequence" gorm:"type:int;not null;default:0"`
	Content          datatypes.JSON `json:"content" gorm:"type:jsonb;not null;default:'{}'"`
	Response		 datatypes.JSON `json:"response" gorm:"type:jsonb;not null;default:'{}'"`
}
