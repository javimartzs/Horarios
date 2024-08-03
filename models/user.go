package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null, unique"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
	Name     string `gorm:"not null"`
}
