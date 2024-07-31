package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
	Name     string `gorm:"not null"`
}
