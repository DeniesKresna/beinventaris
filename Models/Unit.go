package Models

import (
	"gorm.io/gorm"
)

type Unit struct {
	gorm.Model
	Name      string `json:"name"`
	UpdaterID uint

	Updater *User
}

type UnitCreate struct {
	Name      string `form:"name" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

type UnitUpdate struct {
	Name      string `form:"name" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

func (b *Unit) TableName() string {
	return "units"
}
