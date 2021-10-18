package Models

import (
	"gorm.io/gorm"
)

type GoodsType struct {
	gorm.Model
	Name      string `json:"name"`
	Code      string `json:"code"`
	UpdaterID uint

	Updater *User
}

type GoodsTypeCreate struct {
	Name      string `form:"name" validate:"required"`
	Code      string `form:"code" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

type GoodsTypeUpdate struct {
	Name      string `json:"name" validate:"required"`
	Code      string `json:"code" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

func (b *GoodsType) TableName() string {
	return "goods_types"
}
