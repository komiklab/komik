package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
