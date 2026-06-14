package models

import "github.com/google/uuid"

type UserRepresentation struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"uniqueIndex;not null"`
}

func NewUserRepresentation(username string) *UserRepresentation {
	return &UserRepresentation{
		ID:       uuid.New(),
		Username: username,
	}
}
