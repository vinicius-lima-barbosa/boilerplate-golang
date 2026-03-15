package models

import "gorm.io/gorm"

type Type int

const (
	ADMIN Type = iota // 0
	USER              // 1
)

type User struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Email       string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Type        Type   `json:"tipo" gorm:"type:int;not null"`
	Description string `json:"description" gorm:"type:varchar(255)"`
}
