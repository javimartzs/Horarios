package models

import "gorm.io/gorm"

type Worker struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Lastname string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Dni      string `gorm:""`
	Cargo    string `gorm:"not null"`
	Store    string `gorm:"not null"`
	Status   string `gorm:"not null"`
	Prueba   string `gorm:"not null"`
}
