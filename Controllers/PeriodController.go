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

func PeriodIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var periods []Models.Period

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Scopes(FilterModel(search, Models.Period{})), &periods)

	Response.Json(c, 200, p)
}

func PeriodList(c *gin.Context) {
	var periods []Models.Period

	Configs.DB.Find(&periods)
	Response.Json(c, 200, periods)
}

func PeriodShow(c *gin.Context) {
	id := c.Param("id")
	var period Models.Period
	err := Configs.DB.Preload("Updater").First(&period, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.PeriodNotFound)
		return
	}

	Response.Json(c, 200, period)
}

func PeriodActive(c *gin.Context) {
	id := c.Param("id")
	var period Models.Period
	err := Configs.DB.Preload("Updater").Where("active", 1).First(&period, id).Error

	if err != nil {
		Response.Json(c, 200, nil)
		return
	}

	Response.Json(c, 200, period)
}

func PeriodStore(c *gin.Context) {
	SetSessionId(c)
	var period Models.Period
	var periodCreate Models.PeriodCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&periodCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(periodCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", periodCreate.Name).First(&Models.Period{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.PeriodExist)
		return
	}

	periodCreate.UpdaterID = SessionId

	if periodCreate.Active == 1 {
		if err := Configs.DB.Model(Models.Period{}).Where("id > ?", 0).Update("active", 0).Error; err != nil {
			Response.Json(c, 409, Translations.PeriodCreateServerError)
			return
		}
	}

	InjectStruct(&periodCreate, &period)
	if err := Configs.DB.Create(&period).Error; err != nil {
		Response.Json(c, 500, Translations.PeriodCreateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.PeriodCreated)
	}
}

func PeriodUpdate(c *gin.Context) {
	SetSessionId(c)
	var period Models.Period
	var periodUpdate Models.PeriodUpdate
	id := c.Param("id")

	//bind and validate request-------------------------
	if err := c.ShouldBind(&periodUpdate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(periodUpdate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", periodUpdate.Name).Where("id != ?", id).First(&Models.Period{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.PeriodExist)
		return
	}

	//--------------------------------------------------
	err = Configs.DB.First(&period, id).Error
	if err != nil {
		Response.Json(c, 404, Translations.PeriodNotFound)
		return
	}

	periodUpdate.UpdaterID = SessionId

	if periodUpdate.Active == 1 {
		if err := Configs.DB.Model(Models.Period{}).Where("id > ?", 0).Update("active", 0).Error; err != nil {
			Response.Json(c, 409, Translations.PeriodCreateServerError)
			return
		}
	}
	InjectStruct(&periodUpdate, &period)
	if err := Configs.DB.Save(&period).Error; err != nil {
		Response.Json(c, 500, Translations.PeriodUpdateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.PeriodUpdated)
	}
}

func PeriodDestroy(c *gin.Context) {
	id := c.Param("id")
	var period Models.Period
	err := Configs.DB.First(&period, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.PeriodNotFound)
		return
	}

	Configs.DB.Delete(&period)

	Response.Json(c, 200, Translations.PeriodDeleted)
	return
}
