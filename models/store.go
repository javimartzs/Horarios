package models

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	Name   string `gorm:"not null"`
	City   string `gorm:"not null"`
	Phone  string `gorm:"not null"`
	Status string `gorm:"not null"`
}
