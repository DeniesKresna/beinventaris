package Models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	Name              string `json:"name"`
	ImageUrl          string `json:"imageUrl"`
	ProcurementDocUrl string `json:"procurementDocUrl"`
	StatusDocUrl      string `json:"statusDocUrl"`
	Nup               uint   `json:"nup"`
	Year              uint   `json:"year"`
	Quantity          uint   `json:"quantity"`
	Price             uint   `json:"price"`
	UnitID            uint
	GoodsTypeID       uint
	UpdaterID         uint

	Updater    *User
	GoodsType  *GoodsType
	Unit       *Unit
	Histories  []History
	Rooms      []History
	Conditions []History
	Periods    []Period `gorm:"many2many:inventory_period;"`
}

type InventoryCreate struct {
	Name        string    `form:"name" validate:"required"`
	Nup         uint      `form:"nup" validate:"required|int|min:1"`
	Year        uint      `form:"year" validate:"required|int|min:1945|max:2100"`
	Quantity    uint      `form:"quantity" validate:"required|int"`
	Price       uint      `form:"price" validate:"required|int"`
	UnitID      uint      `validate:"required|int"`
	GoodsTypeID uint      `validate:"required|int"`
	HistoryTime time.Time `form:"history_time"`
	Description string    `form:"description"`
	UpdaterID   uint      `validate:"-"`
}

type InventoryUpdate struct {
	Name        string `form:"name" validate:"required"`
	Nup         uint   `form:"nup" validate:"required|int|min:1"`
	Year        uint   `form:"year" validate:"required|int|min:1945|max:2100"`
	Quantity    uint   `form:"quantity" validate:"required|int"`
	Price       uint   `form:"price" validate:"required|int"`
	UnitID      uint   `validate:"required|int"`
	GoodsTypeID uint   `validate:"required|int"`
	UpdaterID   uint   `validate:"-"`
}

func (b *Inventory) TableName() string {
	return "inventories"
}

type InventoryFilterField struct {
	GoodsType uint `form:"goods-type"`
	Unit      uint `form:"unit"`
	Room      uint `form:"room"`
	Condition uint `form:"condition"`
	Period    uint `form:"period"`
}
