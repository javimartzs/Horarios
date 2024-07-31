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
	WeekID   string `gorm:"index:idx_week_worker_interval_day,unique"` // Índice único
	WorkerID uint   `gorm:"index:idx_week_worker_interval_day,unique"`
	Interval string `gorm:"index:idx_week_worker_interval_day,unique"`
	DayIndex int    `gorm:"index:idx_week_worker_interval_day,unique"`
	Color    string
}
