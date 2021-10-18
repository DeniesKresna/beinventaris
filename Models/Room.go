package Models

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name      string `json:"name"`
	UpdaterID uint

	Updater *User
}

type RoomCreate struct {
	Name      string `form:"name" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

type RoomUpdate struct {
	Name      string `form:"name" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

func (b *Room) TableName() string {
	return "rooms"
}
