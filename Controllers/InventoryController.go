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
	"gorm.io/gorm"
)

func InventoryIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var inventories []Models.Inventory

	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Updater").Preload("GoodsType").Preload("Unit").
		Preload("Conditions", func(db *gorm.DB) *gorm.DB {
			return db.Where("histories.entity_type", "condition").Order("histories.created_at DESC").Limit(1)
		}).
		Preload("Rooms", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Room").Where("histories.entity_type", "room").Order("histories.created_at DESC").Limit(1)
		}).Omit("Histories").Scopes(FilterModel(search, Models.Inventory{})), &inventories)

	Response.Json(c, 200, p)
}

func InventoryList(c *gin.Context) {
	var inventories []Models.Inventory

	Configs.DB.Find(&inventories)
	Response.Json(c, 200, inventories)
}

func InventoryShow(c *gin.Context) {
	goodsTypeId, ok := c.GetQuery("goods-type-id")
	if !ok {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}
	nup, ok := c.GetQuery("nup")
	if !ok {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	var inventory Models.Inventory
	err := Configs.DB.Preload("Updater").Preload("GoodsType").Preload("Unit").
		Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("histories.created_at DESC")
		}).Where("goods_type_id", goodsTypeId).Where("nup", nup).First(&inventory).Error

	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	Response.Json(c, 200, inventory)
}

func InventoryStore(c *gin.Context) {
	SetSessionId(c)

	var inventory Models.Inventory
	var inventoryCreate Models.InventoryCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&inventoryCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(inventoryCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------
	err := Configs.DB.Where("goods_type_id = ?", inventoryCreate.GoodsTypeID).Where("nup = ?", inventoryCreate.Nup).First(&Models.Inventory{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.InventoryExist)
		return
	}

	inventoryCreate.UpdaterID = SessionId

	InjectStruct(&inventoryCreate, &inventory)
	if err := Configs.DB.Create(&inventory).Error; err != nil {
		Response.Json(c, 500, Translations.InventoryCreateServerError)
		return
	} else {
		// upload inventory image
		file, err := c.FormFile("image")
		if err != nil {
			Configs.DB.Unscoped().Delete(&inventory)
			Response.Json(c, 500, Translations.InventoryCreateUploadError)
			return
		}
		filename := "inventory-" + strconv.FormatUint(uint64(inventory.ID), 10) + "-" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "-")
		if err := c.SaveUploadedFile(file, Helpers.InventoryPath(filename)); err != nil {
			Configs.DB.Unscoped().Delete(&inventory)
			Response.Json(c, 500, Translations.InventoryCreateUploadError)
			return
		}
		if err := Configs.DB.Model(&inventory).Update("image_url", Helpers.InventoryPath(filename)).Error; err != nil {
			Configs.DB.Unscoped().Delete(&inventory)
			Response.Json(c, 500, Translations.InventoryCreateUploadError)
			return
		}

		var documentsLoop = []map[string]string{
			{"doc": "procurementDoc", "field": "procurement_doc_url"},
			{"doc": "statusDoc", "field": "status_doc_url"},
		}
		// upload inventory documents

		for _, v := range documentsLoop {
			// upload inventory image
			docFile, err := c.FormFile(v["doc"])
			if err != nil {
				continue
			}
			docfilename := "inventory-" + v["doc"] + strconv.FormatUint(uint64(inventory.ID), 10) + "-" + docFile.Filename
			docfilename = strings.ReplaceAll(docfilename, " ", "-")
			if err := c.SaveUploadedFile(docFile, Helpers.InventoryDocumentsPath(docfilename)); err != nil {
				continue
			}
			if err := Configs.DB.Model(&inventory).Update(v["field"], Helpers.InventoryDocumentsPath(docfilename)).Error; err != nil {
				continue
			}
		}

		var historyCreate Models.HistoryCreate

		if err := c.ShouldBind(&historyCreate); err == nil {
			historyCreate.InventoryID = inventory.ID
			historyCreate.EntityType = "room"
			historyCreate.UpdaterID = SessionId

			var history Models.History
			InjectStruct(&historyCreate, &history)
			if err := Configs.DB.Create(&history).Error; err == nil {
				historyFile, err := c.FormFile("historyImage")
				if err == nil {
					filename := "history-" + strconv.FormatUint(uint64(history.ID), 10) + "-" + historyFile.Filename
					filename = strings.ReplaceAll(filename, " ", "-")
					if err := c.SaveUploadedFile(historyFile, Helpers.HistoryPath(filename)); err == nil {
						if err := Configs.DB.Model(&history).Update("image_url", Helpers.HistoryPath(filename)).Error; err != nil {
							Configs.DB.Unscoped().Delete(&history)
						}
					} else {
						Configs.DB.Unscoped().Delete(&history)
					}
				} else {
					Configs.DB.Unscoped().Delete(&history)
				}
			}
		} else {
			Response.Json(c, 500, err)
			return
		}

		Response.Json(c, 200, Translations.InventoryCreated)
	}
}

func InventoryUpdate(c *gin.Context) {
	SetSessionId(c)

	var inventory Models.Inventory
	var inventoryUpdate Models.InventoryUpdate
	id := c.Param("id")

	//bind and validate request-------------------------
	if err := c.ShouldBind(&inventoryUpdate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(inventoryUpdate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------
	err := Configs.DB.Where("goods_type_id = ?", inventoryUpdate.GoodsTypeID).Where("nup = ?", inventoryUpdate.Nup).Where("id != ?", id).First(&Models.Inventory{}).Error
	if err == nil {
		Response.Json(c, 409, Translations.InventoryExist)
		return
	}

	err = Configs.DB.First(&inventory, id).Error
	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	inventoryUpdate.UpdaterID = SessionId

	InjectStruct(&inventoryUpdate, &inventory)
	if err := Configs.DB.Save(&inventory).Error; err != nil {
		Response.Json(c, 500, Translations.InventoryUpdateServerError)
		return
	} else {
		// upload inventory image
		file, err := c.FormFile("image")
		if err == nil {
			filename := "inventory-" + strconv.FormatUint(uint64(inventory.ID), 10) + "-" + file.Filename
			filename = strings.ReplaceAll(filename, " ", "-")
			Helpers.DeleteFile(filename)
			if err := c.SaveUploadedFile(file, Helpers.InventoryPath(filename)); err == nil {
				if err := Configs.DB.Model(&inventory).Update("image_url", Helpers.InventoryPath(filename)).Error; err != nil {

				}
			}
		}

		var documentsLoop = []map[string]string{
			{"doc": "procurementDoc", "field": "procurement_doc_url"},
			{"doc": "statusDoc", "field": "status_doc_url"},
		}
		// upload inventory documents
		for _, v := range documentsLoop {
			// upload inventory image
			docFile, err := c.FormFile(v["doc"])
			if err != nil {
				continue
			}
			docfilename := "inventory-" + v["doc"] + strconv.FormatUint(uint64(inventory.ID), 10) + "-" + docFile.Filename
			docfilename = strings.ReplaceAll(docfilename, " ", "-")
			Helpers.DeleteFile(docfilename)
			if err := c.SaveUploadedFile(docFile, Helpers.InventoryDocumentsPath(docfilename)); err != nil {
				continue
			}
			if err := Configs.DB.Model(&inventory).Update(v["field"], Helpers.InventoryDocumentsPath(docfilename)).Error; err != nil {
				continue
			}
		}

		Response.Json(c, 200, Translations.InventoryUpdated)
	}
}

func InventoryDestroy(c *gin.Context) {
	id := c.Param("id")
	var inventory Models.Inventory
	err := Configs.DB.First(&inventory, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	Configs.DB.Delete(&inventory)

	Response.Json(c, 200, Translations.InventoryDeleted)
	return
}
