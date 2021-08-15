package Models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name      string `json:"name" validate:"required"`
	CreatorID uint

	Creator *User
}

func (b *Role) TableName() string {
	return "roles"
}

type RoleUpdate struct {
	Name      string `validate:"required"`
	CreatorID uint
}
