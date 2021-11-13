package Models

import (
	"gorm.io/gorm"
)

type Survey struct {
	gorm.Model
	EasyUse   int    `json:"easyUse"`
	EasyHelp  int    `json:"easyHelp"`
	Faster    int    `json:"faster"`
	EasyData  int    `json:"easyData"`
	Input     string `json:"input"`
	UpdaterID uint

	Updater *User
}

type SurveyCreate struct {
	EasyUse   int    `json:"easy_use" validate:"required"`
	EasyHelp  int    `json:"easy_help" validate:"required"`
	Faster    int    `json:"faster" validate:"required"`
	EasyData  int    `json:"easy_data" validate:"required"`
	Input     string `json:"input" validate:"required"`
	UpdaterID uint   `validate:"-"`
}

func (b *Survey) TableName() string {
	return "surveys"
}
