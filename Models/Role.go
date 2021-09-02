package Models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name      string `json:"name" validate:"required"`
	UpdaterID uint

	Updater *User
}

func (b *Role) TableName() string {
	return "roles"
}

type RoleUpdate struct {
	Name      string `validate:"required"`
	UpdaterID uint
}
