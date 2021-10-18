package Models

import (
	"gorm.io/gorm"
)

type Condition struct {
	gorm.Model
	Name      string `json:"name"`
	UpdaterID uint

	Updater *User
}

type ConditionCreate struct {
	Name      string `form:"name" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

type ConditionUpdate struct {
	Name      string `form:"name" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

func (b *Condition) TableName() string {
	return "conditions"
}
