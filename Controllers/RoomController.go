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

func RoomIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var rooms []Models.Room

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Scopes(FilterModel(search, Models.Room{})), &rooms)

	Response.Json(c, 200, p)
}

func RoomList(c *gin.Context) {
	var rooms []Models.Room

	Configs.DB.Find(&rooms)
	Response.Json(c, 200, rooms)
}

func RoomShow(c *gin.Context) {
	id := c.Param("id")
	var Room Models.Room
	err := Configs.DB.Preload("Updater").First(&Room, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.RoomNotFound)
		return
	}

	Response.Json(c, 200, Room)
}

func RoomStore(c *gin.Context) {
	SetSessionId(c)
	var Room Models.Room
	var RoomCreate Models.RoomCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&RoomCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(RoomCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", RoomCreate.Name).First(&Models.Room{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.RoomExist)
		return
	}

	RoomCreate.UpdaterID = SessionId

	InjectStruct(&RoomCreate, &Room)
	if err := Configs.DB.Create(&Room).Error; err != nil {
		Response.Json(c, 500, Translations.RoomCreateServerError)
		return
	} else {
		Response.Json(c, 200, Translations.RoomCreated)
	}
}

func RoomDestroy(c *gin.Context) {
	id := c.Param("id")
	var Room Models.Room
	err := Configs.DB.First(&Room, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.RoomNotFound)
		return
	}

	Configs.DB.Delete(&Room)

	Response.Json(c, 200, Translations.RoomDeleted)
	return
}
