package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/komiklab/komik/apidefn"
	"gorm.io/datatypes"
)

const (
	ChannelTypeSlack = "slack"
)

type Channel struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"unique;type:varchar(255);not null"`
	Type      string         `json:"type" gorm:"type:varchar(255);not null"`
	Config    datatypes.JSON `json:"config" gorm:"type:text;not null"`
	CreatedAt int64         `json:"created_at" gorm:"autoCreateTime;not null"`
	UpdatedAt int64         `json:"updated_at" gorm:"autoUpdateTime;not null"`
}

func CreateChannelFromRequest(req apidefn.ChannelRequest) (*Channel, error) {
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, err
	}
	return &Channel{
		Name:   req.Name,
		Type:   string(req.TypeOf),
		Config: datatypes.JSON(payloadBytes),
	}, nil
}

type Channellist struct {
	Channels []Channel `json:"channels"`
}
