package models

import (
	"time"

	"gorm.io/gorm"
)

type Type int

const (
	ADMIN Type = iota // 0
	USER              // 1
)

type User struct {
	gorm.Model
	Name             string     `json:"name" gorm:"type:varchar(100);not null"`
	Email            string     `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash     string     `json:"-" gorm:"type:varchar(255)"`
	RefreshTokenHash *string    `json:"-" gorm:"type:varchar(255)"`
	IsActive         bool       `json:"is_active" gorm:"type:boolean;default:true;not null"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	Type             Type       `json:"tipo" gorm:"type:int;not null"`
	Description      string     `json:"description" gorm:"type:varchar(255)"`
}
