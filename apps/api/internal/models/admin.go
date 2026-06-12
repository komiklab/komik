package models

import (
	"time"
)

type Admin struct {
	ID          uint      `valid:"-" gorm:"primaryKey;autoIncrement"`
	Username    string    `valid:"email,required~Username is required and must be a valid email" gorm:"unique;not null"`
	Password    string    `valid:"length(6|20),required~Password is required  must be between 6 and 20 characters" gorm:"not null"`
	CreatedTime time.Time `valid:"-" gorm:"autoCreateTime"`
	UpdatedTime time.Time `valid:"-" gorm:"autoUpdateTime"`
}

func NewAdminDAO(username, password string) *Admin {
	return &Admin{
		Username: username,
		Password: password,
	}
}
