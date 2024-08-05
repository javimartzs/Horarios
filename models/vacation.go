package models

import "gorm.io/gorm"

type Vacation struct {
	gorm.Model
	WorkerID  uint   `gorm:"not null"`
	Worker    Worker `gorm:"foreignKey:WorkerID"`
	StartDate string `gorm:"not null"`
	EndDate   string `gorm:"not null"`
	Status    string `gorm:"not null"`
}
