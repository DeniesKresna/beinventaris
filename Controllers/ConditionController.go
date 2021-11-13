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

func ConditionIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var conditions []Models.Condition
	var count int64

	Configs.DB.Model(Models.Condition{}).Scopes(FilterModel(search, Models.Condition{})).Count(&count)

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Scopes(FilterModel(search, Models.Condition{})), &conditions, count)

	Response.Json(c, 200, p)
}

func ConditionList(c *gin.Context) {
	var conditions []Models.Condition

	Configs.DB.Find(&conditions)
	Response.Json(c, 200, conditions)
}

func ConditionShow(c *gin.Context) {
	id := c.Param("id")
	var condition Models.Condition
	err := Configs.DB.Preload("Updater").First(&condition, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.ConditionNotFound)
		return
	}

	Response.Json(c, 200, condition)
}

func ConditionStore(c *gin.Context) {
	SetSessionId(c)
	var condition Models.Condition
	var conditionCreate Models.ConditionCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&conditionCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(conditionCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", conditionCreate.Name).First(&Models.Condition{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.ConditionExist)
		return
	}

	conditionCreate.UpdaterID = SessionId

	InjectStruct(&conditionCreate, &condition)
	if err := Configs.DB.Create(&condition).Error; err != nil {
		Response.Json(c, 500, Translations.ConditionCreateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.ConditionCreated)
	}
}

func ConditionDestroy(c *gin.Context) {
	id := c.Param("id")
	var condition Models.Condition
	err := Configs.DB.First(&condition, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.ConditionNotFound)
		return
	}

	Configs.DB.Delete(&condition)

	Response.Json(c, 200, Translations.ConditionDeleted)
	return
}
