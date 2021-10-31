package Controllers

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
	excelize "github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func InventoryPeriodIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	//search := c.DefaultQuery("search", "")
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

	var period Models.Period
	if filtered.Period == 0 {
		if err := Configs.DB.Model(Models.Period{}).Where("active", 1).First(&period).Error; err == nil {
			filtered.Period = period.ID
		} else {
			Response.Json(c, 422, Translations.PeriodNotFound)
			return
		}
	}

	var res Result
	var count int64

	subQueryRoom := Configs.DB.Select("inventory_id,room_id").Where("entity_type = ?", "room").Order("history_time DESC").Table("histories")
	subQueryCondition := Configs.DB.Select("inventory_id,condition_id").Where("entity_type = ?", "condition").Order("history_time DESC").Table("histories")

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
		Where("i.deleted_at is NULL").Scopes(InventoryFilterData(filtered))

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

func InventoryPeriodDelete(c *gin.Context) {
	var period Models.Period
	var inventory Models.Inventory
	inventoryPeriod := struct {
		InventoryId uint `json:"inventory_id"`
		PeriodId    uint `json:"period_id"`
	}{}

	if err := c.ShouldBind(&inventoryPeriod); err != nil {
		Response.Json(c, 422, err)
		return
	}

	if inventoryPeriod.PeriodId == 0 {
		if err := Configs.DB.Model(Models.Period{}).Where("active", 1).First(&period).Error; err != nil {
			Response.Json(c, 422, Translations.PeriodNotFound)
			return
		}
	} else {
		if err := Configs.DB.Model(Models.Period{}).First(&period, inventoryPeriod.PeriodId).Error; err != nil {
			Response.Json(c, 422, Translations.PeriodNotFound)
			return
		}
	}

	if err := Configs.DB.Model(Models.Inventory{}).First(&inventory, inventoryPeriod.InventoryId).Error; err != nil {
		Response.Json(c, 422, Translations.InventoryNotFound)
		return
	}

	Configs.DB.Model(&inventory).Association("Periods").Delete([]Models.Period{period})

	Response.Json(c, 200, Translations.InventoryUpdatePeriodSuccess)
	return
}

func InventoryPeriodExport(c *gin.Context) {
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

	var period Models.Period
	var periodID, ok = c.GetQuery("period_id")
	var periodIDInt int = 0
	if ok {
		periodIDInt, _ = strconv.Atoi(periodID)
	}
	if periodIDInt == 0 {
		if err := Configs.DB.Model(Models.Period{}).Where("active", 1).First(&period).Error; err == nil {
			periodIDInt = int(period.ID)
		}
	}

	subQueryRoom := Configs.DB.Select("inventory_id,room_id").Where("entity_type = ?", "room").Order("history_time DESC").Table("histories")
	subQueryCondition := Configs.DB.Select("inventory_id,condition_id").Where("entity_type = ?", "condition").Order("history_time DESC").Table("histories")

	var query = Configs.DB.Table("inventories as i").Select(`i.id,i.name,i.nup,i.year,
	i.quantity,i.price,un.name as unit_name, gt.id as goods_type_id, gt.name as type_name, gt.code, r.name as room_name, c.name as condition_name
	`).Joins("left join (?) as hr on hr.inventory_id = i.id ", subQueryRoom).
		Joins("left join (?) as hc on hc.inventory_id = i.id ", subQueryCondition).
		Joins("left join rooms as r on r.id = hr.room_id").
		Joins("left join conditions as c on c.id = hc.condition_id").
		Joins("left join units as un on un.id = i.unit_id").
		Joins("left join goods_types as gt on gt.id = i.goods_type_id").
		Joins("left join users as us on us.id = i.updater_id").
		Joins("left join inventory_period as ip on ip.inventory_id = i.id").
		Where("i.deleted_at is NULL").Where("ip.period_id", periodIDInt)

	query.Scan(&inventories)

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

func InventoryFilterData(filterData Models.InventoryFilterField) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filterData.Condition != 0 {
			db = db.Where("c.id", filterData.Condition)
		}
		if filterData.GoodsType != 0 {
			db = db.Where("gt.id", filterData.GoodsType)
		}
		if filterData.Period != 0 {
			db = db.Where("ip.period_id", filterData.Period)
		}
		if filterData.Room != 0 {
			db = db.Where("r.id", filterData.Room)
		}
		if filterData.Unit != 0 {
			db = db.Where("un.id", filterData.Unit)
		}
		return db
	}
}
