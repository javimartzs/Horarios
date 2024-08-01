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

type ScheduleEntry struct {
	gorm.Model
	WeekID   uint   `gorm:"index:idx_week_worker_interval_day,unique"` // Índice único
	WorkerID uint   `gorm:"index:idx_week_worker_interval_day,unique"`
	Interval string `gorm:"index:idx_week_worker_interval_day,unique"`
	DayIndex int    `gorm:"index:idx_week_worker_interval_day,unique"`
	Color    string
}

type WorkerTotal struct {
	gorm.Model
	ID         uint    `gorm:"primaryKey"`
	WorkerID   uint    `gorm:"not null"`
	WeekID     uint    `gorm:"not null"`
	DayIndex   int     `gorm:"not null"`
	TotalHours float64 `gorm:"not null"`
}
