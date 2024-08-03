package models

import "gorm.io/gorm"

type Week struct {
	gorm.Model
	WeekID string `gorm:"not null, unique"`
	Year   int    `gorm:"not null"`
	Week   int    `gorm:"not null"`
	Start  string `gorm:"not null"`
	End    string `gorm:"not null"`
}

type ScheduleEntry struct {
	gorm.Model
	WeekID   int    `gorm:"index:idx_week_worker_interval_day, unique"`
	WorkerID int    `gorm:"index:idx_week_worker_interval_day, unique"`
	Interval string `gorm:"index:idx_week_worker_interval_day, unique"`
	DayIndex int    `gorm:"index:idx_week_worker_interval_day, unique"`
	Color    string
}

type WorkerHours struct {
	gorm.Model
	WorkerID   int     `gorm:"not null"`
	WeekID     int     `gorm:"not null"`
	DayIndex   int     `gorm:"not null"`
	TotalHours float64 `gorm:"not null"`
}
