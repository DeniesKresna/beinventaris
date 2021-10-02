package Controllers

import (
	"strconv"
	"strings"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Helpers"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func HistoryList(c *gin.Context) {
	var histories []Models.History

	Configs.DB.Find(&histories)
	Response.Json(c, 200, histories)
}

func HistoryStore(c *gin.Context) {
	SetSessionId(c)
	var history Models.History
	var historyCreate Models.HistoryCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&historyCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(historyCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------
	historyCreate.UpdaterID = SessionId

	InjectStruct(&historyCreate, &history)
	if err := Configs.DB.Create(&history).Error; err != nil {
		Response.Json(c, 500, Translations.HistoryCreateServerError)
		return
	} else {
		file, err := c.FormFile("image")
		if err != nil {
			Configs.DB.Unscoped().Delete(&history)
			Response.Json(c, 500, Translations.HistoryCreateUploadError)
			return
		}
		filename := "history-" + strconv.FormatUint(uint64(history.ID), 10) + "-" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "-")
		if err := c.SaveUploadedFile(file, Helpers.HistoryPath(filename)); err != nil {
			Configs.DB.Unscoped().Delete(&history)
			Response.Json(c, 500, Translations.HistoryCreateUploadError)
			return
		}
		if err := Configs.DB.Model(&history).Update("image_url", Helpers.HistoryPath(filename)).Error; err != nil {
			Configs.DB.Unscoped().Delete(&history)
			Response.Json(c, 500, Translations.HistoryCreateUploadError)
			return
		}

		Response.Json(c, 200, Translations.HistoryCreated)
	}
}

func HistoryDestroy(c *gin.Context) {
	id := c.Param("id")
	var history Models.History
	err := Configs.DB.First(&history, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.HistoryNotFound)
		return
	}

	Configs.DB.Delete(&history)

	Response.Json(c, 200, Translations.HistoryDeleted)
	return
}
