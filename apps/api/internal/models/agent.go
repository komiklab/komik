package models

import (
	"github.com/asaskevich/govalidator/v12"
	"github.com/google/uuid"
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/utils"
)


//TODO: Do a proper validation[KOM-20]

type Agent struct {
	Id           uuid.UUID                             `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt    int64                                 `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    int64                                 `json:"updatedAt" gorm:"autoUpdateTime"`
	Description  string                                `json:"description" gorm:"null"`
	Name         string                                `json:"name"  valid:"required" gorm:"not null"`
	Endpoint     string                                `json:"endpoint" gorm:"not null"`
	Capabilities []string                              `json:"capabilities" gorm:"type:jsonb;serializer:json;not null"`
	Parameter    []apidefn.AgentCreateRequestParameter `json:"parameter" gorm:"type:jsonb;serializer:json;not null"`
	Tags         *[]string                             `json:"tags" gorm:"type:jsonb;serializer:json;not null"`
}

// type Resources struct {
// 	CPU              string `json:"cpu" gorm:"not null"`
// 	Memory           string `json:"memory" gorm:"not null"`
// 	EphemeralStorage string `json:"ephemeralStorage" gorm:"not null"`
// 	TimeoutSeconds   int    `json:"timeoutSeconds" gorm:"not null"`
// }

// type AgentParameter struct {
// 	Name     string `json:"name" gorm:"not null"`
// 	Type     string `json:"type" gorm:"not null"`
// 	Required bool   `json:"required" gorm:"not null"`
// }

type ListAgents struct {
	Agents []Agent `json:"agents"`
}

func NewAgentFromApiDefn(req apidefn.AgentCreateRequest) (*Agent, error) {
	agent := &Agent{
		Id:           uuid.New(),
		Description:  req.Description,
		Endpoint:     req.Endpoint,
		Capabilities: req.Capabilities,
		Parameter:    req.Parameters,
		Tags:         req.Tags,
		Name:         req.Name,
	}
	valid, err := govalidator.ValidateStruct(agent)
	if !valid {
		return nil, utils.NewBindError("failed to validate agent", err)
	}
	return agent, nil
}
