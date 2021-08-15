package Models

import (
	"gorm.io/gorm"
)

type Academy struct {
	gorm.Model
	Name      string `json:"name"`
	ImageUrl  string `json:"image_url"`
	CreatorID uint

	Creator *User
}

type AcademyCreate struct {
	Name      string `form:"name" validate:"required"`
	CreatorID uint   `validate:"-"`
}

type AcademyUpdate struct {
	Name      string `form:"name" validate:"required"`
	CreatorID uint   `validate:"-"`
}

func (b *Academy) TableName() string {
	return "academies"
}
