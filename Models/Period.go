package Models

import (
	"time"

	"gorm.io/gorm"
)

type Period struct {
	gorm.Model
	Name      string    `json:"name" validate:"required"`
	StartTime time.Time `json:"start_time"`
	Active    int       `json:"active"`
	UpdaterID uint

	Updater *User
}

func (b *Period) TableName() string {
	return "periods"
}

type PeriodCreate struct {
	Name      string    `validate:"required"`
	StartTime time.Time `json:"start_time"`
	Active    int       `json:"active"`
	UpdaterID uint
}

type PeriodUpdate struct {
	Name      string    `validate:"required"`
	StartTime time.Time `json:"start_time"`
	Active    int       `json:"active"`
	UpdaterID uint
}
