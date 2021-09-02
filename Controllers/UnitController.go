package Controllers

import (
	"fmt"
	"strconv"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func UnitIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var units []Models.Unit

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Scopes(FilterModel(search, Models.Unit{})), &units)

	Response.Json(c, 200, p)
}

func UnitList(c *gin.Context) {
	var units []Models.Unit

	Configs.DB.Find(&units)
	Response.Json(c, 200, units)
}

func UnitShow(c *gin.Context) {
	id := c.Param("id")
	var Unit Models.Unit
	err := Configs.DB.Preload("Updater").First(&Unit, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.UnitNotFound)
		return
	}

	Response.Json(c, 200, Unit)
}

func UnitStore(c *gin.Context) {
	SetSessionId(c)
	var Unit Models.Unit
	var UnitCreate Models.UnitCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&UnitCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	fmt.Print(UnitCreate)
	v := validate.Struct(UnitCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", UnitCreate.Name).First(&Models.Unit{}).Error
	if err == nil {
		Response.Json(c, 500, Translations.UnitExist)
		return
	}

	UnitCreate.UpdaterID = SessionId

	InjectStruct(&UnitCreate, &Unit)
	if err := Configs.DB.Create(&Unit).Error; err != nil {
		Response.Json(c, 500, Translations.UnitCreateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.UnitCreated)
	}
}

func UnitDestroy(c *gin.Context) {
	id := c.Param("id")
	var Unit Models.Unit
	err := Configs.DB.First(&Unit, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.UnitNotFound)
		return
	}

	Configs.DB.Delete(&Unit)

	Response.Json(c, 200, Translations.UnitDeleted)
	return
}
