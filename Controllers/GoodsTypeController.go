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

func GoodsTypeIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var goodsTypes []Models.GoodsType
	var count int64

	Configs.DB.Model(Models.GoodsType{}).Scopes(FilterModel(search, Models.GoodsType{})).Count(&count)

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Scopes(FilterModel(search, Models.GoodsType{})), &goodsTypes, count)

	Response.Json(c, 200, p)
}

func GoodsTypeList(c *gin.Context) {
	var goodsTypes []Models.GoodsType

	Configs.DB.Find(&goodsTypes)
	Response.Json(c, 200, goodsTypes)
}

func GoodsTypeShow(c *gin.Context) {
	id := c.Param("id")
	var goodsType Models.GoodsType
	err := Configs.DB.Preload("Updater").First(&goodsType, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.GoodsTypeNotFound)
		return
	}

	Response.Json(c, 200, goodsType)
}

func GoodsTypeStore(c *gin.Context) {
	SetSessionId(c)
	var goodsType Models.GoodsType
	var goodsTypeCreate Models.GoodsTypeCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&goodsTypeCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(goodsTypeCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", goodsTypeCreate.Name).Or("code = ?", goodsTypeCreate.Code).First(&Models.GoodsType{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.GoodsTypeExist)
		return
	}

	goodsTypeCreate.UpdaterID = SessionId

	InjectStruct(&goodsTypeCreate, &goodsType)
	if err := Configs.DB.Create(&goodsType).Error; err != nil {
		Response.Json(c, 500, Translations.GoodsTypeCreateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.GoodsTypeCreated)
	}
}

func GoodsTypeUpdate(c *gin.Context) {
	SetSessionId(c)
	var goodsType Models.GoodsType
	var goodsTypeUpdate Models.GoodsTypeUpdate
	id := c.Param("id")

	//bind and validate request-------------------------
	if err := c.ShouldBindJSON(&goodsTypeUpdate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(goodsTypeUpdate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where(Configs.DB.Where("name = ?", goodsTypeUpdate.Name).Or("code = ?", goodsTypeUpdate.Code)).Where("id != ?", id).First(&Models.GoodsType{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.GoodsTypeExist)
		return
	}

	err = Configs.DB.First(&goodsType, id).Error
	if err != nil {
		Response.Json(c, 404, Translations.GoodsTypeNotFound)
		return
	}

	goodsTypeUpdate.UpdaterID = SessionId
	InjectStruct(&goodsTypeUpdate, &goodsType)
	if err := Configs.DB.Save(&goodsType).Error; err != nil {
		Response.Json(c, 500, Translations.GoodsTypeUpdateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.GoodsTypeUpdated)
	}
}

func GoodsTypeDestroy(c *gin.Context) {
	id := c.Param("id")
	var goodsType Models.GoodsType
	err := Configs.DB.First(&goodsType, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.GoodsTypeNotFound)
		return
	}

	Configs.DB.Delete(&goodsType)

	Response.Json(c, 200, Translations.GoodsTypeDeleted)
	return
}
