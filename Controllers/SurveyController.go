package Controllers

import (
	"strconv"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func SurveyIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var surveys []Models.Survey
	var count int64

	Configs.DB.Model(Models.Survey{}).Scopes(FilterModel(search, Models.Survey{})).Count(&count)

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Scopes(FilterModel(search, Models.Survey{})), &surveys, count)

	Response.Json(c, 200, p)
}

func SurveyList(c *gin.Context) {
	var surveys []Models.Survey

	Configs.DB.Find(&surveys)
	Response.Json(c, 200, surveys)
}

func SurveyShow(c *gin.Context) {
	id := c.Param("id")
	var Survey Models.Survey
	err := Configs.DB.Preload("Updater").First(&Survey, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.SurveyNotFound)
		return
	}

	Response.Json(c, 200, Survey)
}

func SurveyMe(c *gin.Context) {
	SetSessionId(c)
	var Survey Models.Survey
	err := Configs.DB.Preload("Updater").Where("updater_id", SessionId).First(&Survey).Error

	if err != nil {
		Response.Json(c, 404, "Kamu belum mengisi Kuesioner")
		return
	}

	Response.Json(c, 200, Survey)
}

func SurveyStore(c *gin.Context) {
	SetSessionId(c)
	var Survey Models.Survey
	var SurveyCreate Models.SurveyCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&SurveyCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(SurveyCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------
	SurveyCreate.UpdaterID = SessionId
	InjectStruct(&SurveyCreate, &Survey)

	if err := Configs.DB.Where("updater_id", SessionId).First(&Models.Survey{}).Error; err != nil {
		if err2 := Configs.DB.Create(&Survey).Error; err2 != nil {
			Response.Json(c, 500, Translations.SurveyCreateServerError)
			return
		}
	} else {
		if err2 := Configs.DB.Save(&Survey).Error; err2 != nil {
			Response.Json(c, 500, Translations.SurveyUpdateServerError)
			return
		}
	}

	Response.Json(c, 200, Translations.SurveyUpdated)
	return
}

func SurveyDestroy(c *gin.Context) {
	id := c.Param("id")
	var Survey Models.Survey
	err := Configs.DB.First(&Survey, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.SurveyNotFound)
		return
	}

	Configs.DB.Delete(&Survey)

	Response.Json(c, 200, Translations.SurveyDeleted)
	return
}
