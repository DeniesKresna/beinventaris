package Models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"-" validate:"required"`
	Email    string `json:"email" validate:"required|email"`
	Phone    string `json:"phone" validate:"required"`
	RoleID   uint

	Role *Role
}

func (b *User) TableName() string {
	return "users"
}

type UserUpdate struct {
	RoleId uint8 `validate:"int"`
}

type UserLogin struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
