package models

import (
	"github.com/asaskevich/govalidator/v12"
	"github.com/google/uuid"
	"github.com/komiklab/komik/apidefn"
	"github.com/komiklab/komik/internal/utils"
)

const (
	ContainerRuntime = "container"
	ApiVersionV1     = "komik/v1"
)

//TODO: Do a proper validation[KOM-20]

type Agent struct {
	Id              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt       int64     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       int64     `json:"updatedAt" gorm:"autoUpdateTime"`
	ApiVersion      string    `json:"apiVersion" gorm:"not null"`
	Description     *string   `json:"description" gorm:"null"`
	Image           string    `json:"image" valid:"required" gorm:"not null"`
	ImagePullSecret *string   `json:"imagePullSecret" gorm:"null"`
	Runtime         string    `json:"runtime" gorm:"not null"`
	Annotations     *[]KVPair `json:"annotations" gorm:"serializer:json;not null"`
	Env             *[]KVPair `json:"env" gorm:"serializer:json;not null"`
	Secrets         *[]KVPair `json:"secrets" gorm:"serializer:json;not null"`
	Name            string    `json:"name"  valid:"required" gorm:"not null"`
	*apidefn.AgentCreateRequestResources
}

// type Resources struct {
// 	CPU              string `json:"cpu" gorm:"not null"`
// 	Memory           string `json:"memory" gorm:"not null"`
// 	EphemeralStorage string `json:"ephemeralStorage" gorm:"not null"`
// 	TimeoutSeconds   int    `json:"timeoutSeconds" gorm:"not null"`
// }

type ListAgents struct {
	Agents []Agent `json:"agents"`
}

type KVPair struct {
	Name  string `json:"name" gorm:"not null"`
	Value string `json:"value" gorm:"not null"`
}

func NewKVPairArray(arr *[]struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}) *[]KVPair {
	kvPairs := make([]KVPair, 0)
	if arr != nil {
		kvPairs = make([]KVPair, 0, len(*arr))
		for _, v := range *arr {
			kvPairs = append(kvPairs, KVPair{
				Name:  v.Name,
				Value: v.Value,
			})
		}
	}
	return &kvPairs
}

func NewAgentFromApiDefn(req apidefn.AgentCreateRequest) (*Agent, error) {
	agent := &Agent{
		Id:                          uuid.New(),
		ApiVersion:                  ApiVersionV1,
		Description:                 req.Description,
		Image:                       req.Image,
		ImagePullSecret:             req.ImagePullSecret,
		Runtime:                     ContainerRuntime,
		Annotations:                 NewKVPairArray(req.Annotations),
		Env:                         NewKVPairArray(req.Env),
		Secrets:                     NewKVPairArray(req.Secrets),
		Name:                        req.Name,
		AgentCreateRequestResources: req.Resources,
	}
	valid, err := govalidator.ValidateStruct(agent)
	if !valid {
		return nil, utils.NewBindError("failed to validate agent", err)
	}
	return agent, nil
}
