package Models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	Name        string `json:"name"`
	UnitID      uint
	GoodsTypeID uint
	UpdaterID   uint

	Updater *User
}

type InventoryCreate struct {
	Name        string `form:"name" validate:"required"`
	UnitID      uint   `validate:"-"`
	GoodsTypeID uint   `validate:"-"`
	UpdaterID   uint   `validate:"-"`
}

type InventoryUpdate struct {
	Name        string `form:"name" validate:"required"`
	UnitID      uint   `validate:"-"`
	GoodsTypeID uint   `validate:"-"`
	UpdaterID   uint   `validate:"-"`
}

func (b *Inventory) TableName() string {
	return "inventories"
}
