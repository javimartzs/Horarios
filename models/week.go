package models

import "gorm.io/gorm"

type Week struct {
	gorm.Model
	Year   int    `gorm:"not null"`
	Week   int    `gorm:"not null"`
	Start  string `gorm:"not null"`
	End    string `gorm:"not null"`
	WeekID string `gorm:"not null, unique"`
}
