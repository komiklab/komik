package models

import (
	"time"
)

type Admin struct {
	ID          uint      `valid:"-" gorm:"primaryKey;autoIncrement"`
	Username    string    `valid:"email,required" gorm:"unique;not null"`
	Password    string    `valid:"min=8,max=72,required" gorm:"not null"`
	CreatedTime time.Time `valid:"-" gorm:"autoCreateTime"`
	UpdatedTime time.Time `valid:"-" gorm:"autoUpdateTime"`
}
