package models

import "gorm.io/gorm"

type Worker struct {
	gorm.Model
	Name           string `gorm:"not null"`
	Lastname       string `gorm:"not null"`
	Email          string `gorm:"not null"`
	Identification string `gorm:"not null"`
	Cargo          string `gorm:"not null"`
	Store          string `gorm:"not null"`
	Status         string `gorm:"not null"`
	PeriodoPrueba  string `gorm:"not null"`
}
