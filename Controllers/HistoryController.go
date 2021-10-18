package Controllers

import (
	"database/sql"
	"math"
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

func HistoryIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	//search := c.DefaultQuery("search", "")
	var histories []Models.HistoryIndexData
	var res Result
	var count int64
	inventoryId := c.Param("id")

	Configs.DB.Model(Models.History{}).Where("inventory_id", inventoryId).Count(&count)
	Configs.DB.Table("histories as h").Select(`h.id, h.inventory_id, h.entity_type, 
		h.history_time, h.room_id, h.condition_id, h.description, h.image_url, u.name as updater_name,
		h.created_at, h.updated_at,
		CASE 
			WHEN h.entity_type = 'room' THEN r.name
			ELSE c.name
		END AS entity_name
	`).Joins("join users as u on u.id = h.updater_id").
		Joins("left join rooms as r on r.id = h.room_id").
		Joins("left join conditions as c on c.id = h.condition_id").
		Where("h.inventory_id", inventoryId).
		Where("h.deleted_at is NULL").Order("h.history_time DESC").Offset(pageSize * (page - 1)).Limit(pageSize).Scan(&histories)

	res.CurrentPage = page
	res.Data = histories
	res.From = (pageSize * (page - 1)) + 1
	res.LastPage = int(math.Ceil(float64(count) / float64(pageSize)))
	res.PerPage = pageSize
	if res.CurrentPage < res.LastPage {
		res.To = int64(pageSize * page)
	} else {
		res.To = count
	}
	res.Total = count

	Response.Json(c, 200, res)
}

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
	if historyCreate.EntityType == "room" {
		historyCreate.RoomID = sql.NullInt32{Int32: int32(historyCreate.EntityID), Valid: true}
		historyCreate.ConditionID = sql.NullInt32{Valid: false}
	} else {
		historyCreate.ConditionID = sql.NullInt32{Int32: int32(historyCreate.EntityID), Valid: true}
		historyCreate.RoomID = sql.NullInt32{Valid: false}
	}

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

func HistoryUpdate(c *gin.Context) {
	SetSessionId(c)
	var history Models.History
	var historyUpdate Models.HistoryUpdate
	id := c.Param("id")

	//bind and validate request-------------------------
	if err := c.ShouldBind(&historyUpdate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(historyUpdate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------
	err := Configs.DB.First(&history, id).Error
	if err != nil {
		Response.Json(c, 404, Translations.HistoryNotFound)
		return
	}

	historyUpdate.UpdaterID = SessionId
	if historyUpdate.EntityType == "room" {
		historyUpdate.RoomID = sql.NullInt32{Int32: int32(historyUpdate.EntityID), Valid: true}
		historyUpdate.ConditionID = sql.NullInt32{Valid: false}
	} else {
		historyUpdate.ConditionID = sql.NullInt32{Int32: int32(historyUpdate.EntityID), Valid: true}
		historyUpdate.RoomID = sql.NullInt32{Valid: false}
	}

	InjectStruct(&historyUpdate, &history)
	if err := Configs.DB.Save(&history).Error; err != nil {
		Response.Json(c, 500, Translations.HistoryUpdateServerError)
		return
	} else {
		file, err := c.FormFile("image")
		if err == nil {
			filename := "history-" + strconv.FormatUint(uint64(history.ID), 10) + "-" + file.Filename
			filename = strings.ReplaceAll(filename, " ", "-")
			_ = Helpers.InventoryDocumentsPath(filename)
			if err := c.SaveUploadedFile(file, Helpers.HistoryPath(filename)); err == nil {
				if err := Configs.DB.Model(&history).Update("image_url", Helpers.HistoryPath(filename)).Error; err == nil {
				}
			}
		}

		Response.Json(c, 200, Translations.HistoryUpdated)
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
