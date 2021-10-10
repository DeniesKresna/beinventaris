package Models

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	InventoryID uint      `json:"inventoryID"`
	EntityType  string    `json:"entity_type"`
	RoomID      uint      `json:"RoomID"`
	ConditionID uint      `json:"ConditionID"`
	HistoryTime time.Time `json:"history_time"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	UpdaterID   uint

	Updater   *User
	Inventory *Inventory
	Room      *Room
	Condition *Condition
}

type HistoryIndexData struct {
	gorm.Model
	InventoryID uint      `json:"InventoryID"`
	EntityType  string    `json:"entity_type"`
	EntityName  string    `json:"entity_name"`
	RoomID      uint      `json:"RoomID"`
	ConditionID uint      `json:"ConditionID"`
	HistoryTime time.Time `json:"history_time"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	UpdaterName string    `json:"updater_name"`
}

type HistoryCreate struct {
	InventoryID uint   `form:"InventoryID" validate:"required|int"`
	EntityType  string `form:"entity_type" validate:"in:room,condition"`
	EntityID    uint   `form:"EntityID" validate:"required|int"`
	RoomID      uint
	ConditionID uint
	HistoryTime time.Time `form:"history_time" validate:"required"`
	Description string    `form:"description" validate:"required"`
	UpdaterID   uint
}

type HistoryUpdate struct {
	InventoryID uint   `form:"InventoryID" validate:"required|int"`
	EntityType  string `form:"entity_type" validate:"in:room,condition"`
	EntityID    uint   `form:"EntityID" validate:"required|int"`
	RoomID      uint
	ConditionID uint
	HistoryTime time.Time `form:"history_time" validate:"required"`
	Description string    `form:"description" validate:"required"`
	UpdaterID   uint
}

func (b *History) TableName() string {
	return "histories"
}
