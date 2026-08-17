package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/komiklab/komik/internal/utils"
)

type Hooks struct {
	Id        uuid.UUID `json:"id" valid:"-" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" valid:"required~Name is required" gorm:"not null;unique"`
	CreatedAt time.Time `json:"createdAt" valid:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" valid:"-" gorm:"autoUpdateTime"`
	Secret    string    `json:"secretKey" valid:"required~Secret is required" gorm:"not null"`
}

func NewHook(name string) *Hooks {
	return &Hooks{
		Id:     uuid.New(),
		Name:   name,
		Secret: utils.Generate32bitSecret(),
	}
}

type ListHooks struct {
	Hooks []Hooks `json:"hooks"`
}
