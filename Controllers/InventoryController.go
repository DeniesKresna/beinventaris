package Controllers

import (
	"database/sql"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Helpers"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	excelize "github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func InventoryIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var inventories = []struct {
		ID                uint      `json:"ID"`
		UpdatedAt         time.Time `json:"updated_at"`
		CreatedAt         time.Time `json:"created_at"`
		Name              string    `json:"name"`
		ImageUrl          string    `json:"imageUrl"`
		ProcurementDocUrl string    `json:"procurementDocUrl"`
		StatusDocUrl      string    `json:"statusDocUrl"`
		Nup               uint      `json:"nup"`
		Year              uint      `json:"year"`
		Quantity          uint      `json:"quantity"`
		Price             uint      `json:"price"`
		UnitName          string    `json:"unit_name"`
		ConditionName     string    `json:"condition_name"`
		RoomName          string    `json:"room_name"`
		GoodsTypeID       string    `json:"goods_type_id"`
		GoodsTypeName     string    `json:"goods_type_name"`
		GoodsTypeCode     string    `json:"goods_type_code"`
		UpdaterName       string    `json:"user_name"`
	}{}

	var filtered Models.InventoryFilterField

	if err := c.Bind(&filtered); err != nil {
		Response.Json(c, 422, Translations.InventoryFilteredNotFound)
		return
	}

	var res Result
	var count int64

	subQueryRoom := Configs.DB.Select("inventory_id,room_id").Where("entity_type = ?", "room").Order("history_time DESC").Limit(1).Table("histories")
	subQueryCondition := Configs.DB.Select("inventory_id,condition_id").Where("entity_type = ?", "condition").Order("history_time DESC").Limit(1).Table("histories")

	var query = Configs.DB.Table("inventories as i").Select(`i.id,i.name,i.nup,i.year,i.deleted_at,i.updated_at, i.created_at, i.procurement_doc_url, i.status_doc_url, i.image_url,
	i.quantity,i.price,un.name as unit_name, gt.id as goods_type_id, gt.name as goods_type_name, gt.code as goods_type_code, hr.room_id, hc.condition_id, r.name as room_name, c.name as condition_name, us.name as updater_name
	`).Joins("left join (?) as hr on hr.inventory_id = i.id ", subQueryRoom).
		Joins("left join (?) as hc on hc.inventory_id = i.id ", subQueryCondition).
		Joins("left join rooms as r on r.id = hr.room_id").
		Joins("left join conditions as c on c.id = hc.condition_id").
		Joins("left join units as un on un.id = i.unit_id").
		Joins("left join goods_types as gt on gt.id = i.goods_type_id").
		Joins("left join users as us on us.id = i.updater_id").
		Joins("left join inventory_period as ip on ip.inventory_id = i.id").
		Where("i.deleted_at is NULL").Where("i.name like ?", "%"+search+"%").Group("i.id").Order("i.updated_at DESC").Scopes(InventoryFilterData(filtered))

	query.Count(&count)
	query.Offset(pageSize * (page - 1)).Limit(pageSize).Scan(&inventories)

	res.CurrentPage = page
	res.Data = inventories
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

func InventoryList(c *gin.Context) {
	var inventories []Models.Inventory

	Configs.DB.Find(&inventories)
	Response.Json(c, 200, inventories)
}

func InventoryShow(c *gin.Context) {
	filterCode := struct {
		GoodsTypeId uint
		Nup         uint
	}{}

	//bind and validate request-------------------------
	if err := c.ShouldBind(&filterCode); err != nil {
		Response.Json(c, 422, err)
		return
	}

	var inventory Models.Inventory
	err := Configs.DB.Preload("Updater").Preload("GoodsType").Preload("Unit").
		Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("histories.created_at DESC")
		}).Where("goods_type_id", filterCode.GoodsTypeId).Where("nup", filterCode.Nup).First(&inventory).Error

	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	Response.Json(c, 200, inventory)
}

func InventoryCodeNameShow(c *gin.Context) {
	filterCode := struct {
		Code string
		Nup  uint
	}{}

	//bind and validate request-------------------------
	if err := c.ShouldBind(&filterCode); err != nil {
		Response.Json(c, 422, err)
		return
	}

	var goodsType Models.GoodsType
	if err := Configs.DB.Where("code", filterCode.Code).First(&goodsType).Error; err != nil {
		Response.Json(c, 404, Translations.GoodsTypeNotFound)
	}

	var inventory Models.Inventory
	err := Configs.DB.Where("goods_type_id", goodsType.ID).Where("nup", filterCode.Nup).First(&inventory).Error

	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	Response.Json(c, 200, inventory)
}

func InventoryShowDetail(c *gin.Context) {
	id := c.Param("id")

	var inventory struct {
		gorm.Model
		Name              string `json:"name"`
		ImageUrl          string `json:"imageUrl"`
		ProcurementDocUrl string `json:"procurementDocUrl"`
		StatusDocUrl      string `json:"statusDocUrl"`
		Nup               uint   `json:"nup"`
		Year              uint   `json:"year"`
		Quantity          uint   `json:"quantity"`
		Price             uint   `json:"price"`
		UnitID            uint
		GoodsTypeID       uint
		UpdaterID         uint

		Updater   *Models.User
		GoodsType *Models.GoodsType
		Unit      *Models.Unit
		Histories []Models.History
		Room      Models.Room
		Condition Models.Condition
		Periods   []Models.Period
	}
	err := Configs.DB.Model(Models.Inventory{}).Preload("Updater").Preload("GoodsType").Preload("Unit").
		Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Updater").
				Preload("Condition").
				Preload("Room")
		}).Where("id", id).First(&inventory).Error

	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	Configs.DB.Model(Models.Room{}).Joins("join histories h on h.room_id = rooms.id").Where("h.inventory_id", inventory.ID).Where("h.entity_type", "room").Order("h.history_time DESC").First(&inventory.Room)
	Configs.DB.Model(Models.Condition{}).Joins("join histories h on h.condition_id = conditions.id").Where("h.inventory_id", inventory.ID).Where("h.entity_type", "condition").Order("h.history_time DESC").First(&inventory.Condition)
	Configs.DB.Model(Models.Period{}).Joins("join inventory_period ip on ip.period_id = periods.id").Where("ip.inventory_id", inventory.ID).Scan(&inventory.Periods)

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
		if err == nil {
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
			/*
				if err := c.SaveUploadedFile(docFile, Helpers.InventoryDocumentsPath(docfilename)); err != nil {
					continue
				}*/
			src, _ := file.Open()
			defer src.Close()

			dst, _ := os.Create(Helpers.InventoryDocumentsPath(docfilename))
			defer dst.Close()

			io.Copy(dst, src)

			if err := Configs.DB.Model(&inventory).Update(v["field"], Helpers.InventoryDocumentsPath(docfilename)).Error; err != nil {
				continue
			}
		}

		var historyCreate Models.HistoryCreate

		if err := c.ShouldBind(&historyCreate); err == nil {
			historyCreate.InventoryID = inventory.ID
			historyCreate.EntityType = "room"
			historyCreate.RoomID = sql.NullInt32{Int32: int32(historyCreate.EntityID), Valid: true}
			historyCreate.ConditionID = sql.NullInt32{Valid: false}
			historyCreate.UpdaterID = SessionId

			var history Models.History
			InjectStruct(&historyCreate, &history)
			if err := Configs.DB.Create(&history).Error; err == nil { /*
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
					}*/
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
			_ = Helpers.DeleteFile(Helpers.InventoryPath(filename))
			if err := c.SaveUploadedFile(file, Helpers.InventoryPath(filename)); err == nil {
				if err := Configs.DB.Model(&inventory).Update("image_url", Helpers.InventoryPath(filename)).Error; err != nil {

				}
			}
		} else {
			_ = Helpers.DeleteFile(inventory.ImageUrl)
			if err := Configs.DB.Model(&inventory).Update("image_url", "").Error; err != nil {

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

func InventoryExport(c *gin.Context) {
	f := excelize.NewFile()
	const sheet = "Rekap"
	index := f.NewSheet(sheet)
	var columns = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
	var data = [][]interface{}{
		{"No", "Kode Barang", "Nama Barang", "Tahun Perolehan", "NUP", "Merk/Type", "Satuan", "Kuantitas", "Harga Satuan Barang", "Harga Barang (Rp)", "Kondisi", "Kondisi", "Kondisi", "Ruangan"},
		{"", "", "", "", "", "", "", "", "", "", "B", "RR", "R", ""},
	}

	for i := 0; i < len(data); i++ {
		for j, val := range data[i] {
			f.SetCellValue(sheet, columns[j]+strconv.Itoa(i+1), val)
		}
	}

	_ = f.MergeCell(sheet, "K1", "M1")

	for j := 0; j < len(columns); j++ {
		if columns[j] == "K" || columns[j] == "L" || columns[j] == "M" {
		} else {
			_ = f.MergeCell(sheet, columns[j]+"1", columns[j]+"2")
		}
	}

	var inventories = []struct {
		Id            uint
		Name          string
		Nup           int
		Year          int
		Quantity      int
		Price         int
		UnitName      string
		TypeName      string
		Code          string
		RoomName      string
		ConditionName string
	}{}

	Configs.DB.Raw(`SELECT i.id, i.name, i.nup, i.year, i.quantity, i.price, u.name as unit_name, g.name as type_name, g.code,
		(SELECT r.name FROM histories h INNER JOIN rooms r ON r.id = h.room_id order by h.history_time DESC LIMIT 1) as room_name,
		(SELECT c.name FROM histories h INNER JOIN conditions c ON c.id = h.condition_id order by h.history_time DESC LIMIT 1) as condition_name
		FROM inventories i LEFT JOIN goods_types g ON g.id = i.goods_type_id LEFT JOIN units u ON u.id = i.unit_id
	`).Scan(&inventories)

	for i, val := range inventories {
		rows := strconv.Itoa(i + 3)
		f.SetCellValue(sheet, "A"+rows, val.Id)
		f.SetCellValue(sheet, "B"+rows, val.Code)
		f.SetCellValue(sheet, "C"+rows, val.Name)
		f.SetCellValue(sheet, "D"+rows, val.Year)
		f.SetCellValue(sheet, "E"+rows, val.Nup)
		f.SetCellValue(sheet, "F"+rows, val.TypeName)
		f.SetCellValue(sheet, "G"+rows, val.UnitName)
		f.SetCellValue(sheet, "H"+rows, val.Quantity)
		f.SetCellValue(sheet, "I"+rows, val.Price)
		f.SetCellValue(sheet, "J"+rows, val.Price*val.Quantity)
		f.SetCellValue(sheet, "K"+rows, "")
		f.SetCellValue(sheet, "L"+rows, "")
		f.SetCellValue(sheet, "M"+rows, "")
		var conditionName = strings.ToLower(val.ConditionName)
		if conditionName == "baik" {
			f.SetCellValue(sheet, "K"+rows, 1)
		} else if conditionName == "rusak ringan" {
			f.SetCellValue(sheet, "L"+rows, 1)
		} else {
			f.SetCellValue(sheet, "M"+rows, 1)
		}
		f.SetCellValue(sheet, "N"+rows, val.RoomName)
	}

	f.SetActiveSheet(index)

	currentTime := time.Now()
	crtFormated := currentTime.Format("02012006")

	buf, _ := f.WriteToBuffer()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+crtFormated+".xlsx")
	c.Data(200, "application/octet-stream", buf.Bytes())
}

func InventoryUpdatePeriod(c *gin.Context) {
	var period Models.Period
	var inventory Models.Inventory
	inventoryPeriod := struct {
		InventoryId uint `json:"inventory_id"`
		PeriodId    uint `json:"period_id"`
		AddMode     bool `json:"add_mode"`
	}{}

	if err := c.ShouldBind(&inventoryPeriod); err != nil {
		Response.Json(c, 422, err)
		return
	}

	if err := Configs.DB.Model(Models.Period{}).First(&period, inventoryPeriod.PeriodId).Error; err != nil {
		Response.Json(c, 422, Translations.PeriodNotFound)
		return
	}

	if err := Configs.DB.Model(Models.Inventory{}).First(&inventory, inventoryPeriod.InventoryId).Error; err != nil {
		Response.Json(c, 422, Translations.InventoryNotFound)
		return
	}

	if inventoryPeriod.AddMode {
		Configs.DB.Model(&inventory).Association("Periods").Append([]Models.Period{period})
	} else {
		Configs.DB.Model(&inventory).Association("Periods").Delete([]Models.Period{period})
	}

	Response.Json(c, 200, Translations.InventoryUpdatePeriodSuccess)
	return
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

func InventoryDownloadDoc(c *gin.Context) {
	inventoryId := c.DefaultQuery("inventoryId", "")
	docType := c.DefaultQuery("docType", "")
	if inventoryId == "" || !(docType == "status" || docType == "procurement") {
		Response.Json(c, 404, Translations.InventoryDocumentNotFound)
		return
	}
	var inventory Models.Inventory

	err := Configs.DB.First(&inventory, inventoryId).Error
	if err != nil {
		Response.Json(c, 404, Translations.InventoryNotFound)
		return
	}

	var url string
	if docType == "procurement" {
		if inventory.ProcurementDocUrl == "" {
			Response.Json(c, 404, Translations.InventoryDocumentNotFound)
			return
		}
		url = inventory.ProcurementDocUrl
	}
	if docType == "status" {
		if inventory.StatusDocUrl == "" {
			Response.Json(c, 404, Translations.InventoryDocumentNotFound)
			return
		}
		url = inventory.StatusDocUrl
	}
	file, err := os.Open("./" + url)
	if err != nil {
		Response.Json(c, 404, Translations.InventoryDocumentNotFound)
		return
	}
	defer file.Close()
	c.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		Response.Json(c, 404, Translations.InventoryDocumentNotFound)
		return
	}
	return
}

func InventoryDownloadDocuments(c *gin.Context) {
	inventoryId := c.DefaultQuery("inventoryId", "")
	docType := c.DefaultQuery("docType", "")
	if inventoryId == "" || !(docType == "status" || docType == "procurement") {
		return
	}
	var inventory Models.Inventory

	err := Configs.DB.First(&inventory, inventoryId).Error
	if err != nil {
		return
	}

	var s = make([]string, 0)
	var url string
	if docType == "procurement" {
		if inventory.ProcurementDocUrl == "" {
			return
		}
		s = strings.Split(inventory.ProcurementDocUrl, "/")
		url = inventory.ProcurementDocUrl
	}
	if docType == "status" {
		if inventory.StatusDocUrl == "" {
			return
		}
		s = strings.Split(inventory.StatusDocUrl, "/")
		url = inventory.StatusDocUrl
	}
	ext := s[len(s)-1]
	c.Header("Content-Description", "Document File Download")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+ext)
	c.Header("Content-Type", "application/octet-stream")

	c.File("./" + url)
	return
}
