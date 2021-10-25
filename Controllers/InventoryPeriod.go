package Controllers

import (
	"math"
	"strconv"
	"time"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
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
	/*
		p, _ := (&PConfig{
			Page:    page,
			PerPage: pageSize,
			Path:    c.FullPath(),
			Sort:    "id desc",
		}).Paginate(Configs.DB.Preload("Updater").Preload("GoodsType").Preload("Unit").
			Preload("Conditions", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Condition").Where("histories.inventory_id", "inventories.id").Where("histories.entity_type", "condition").Order("histories.created_at DESC").Limit(1)
			}).
			Preload("Rooms", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Room").Where("histories.inventory_id", "inventories.id").Where("histories.entity_type", "room").Order("histories.created_at DESC").Limit(1)
			}).Omit("Histories").Scopes(FilterModel(search, Models.Inventory{})), &inventories)
	*/
	var res Result
	var count int64

	Configs.DB.Model(Models.Inventory{}).Count(&count)
	subQueryRoom := Configs.DB.Select("inventory_id,room_id").Where("entity_type = ?", "room").Order("history_time DESC").Table("histories")
	subQueryCondition := Configs.DB.Select("inventory_id,condition_id").Where("entity_type = ?", "condition").Order("history_time DESC").Table("histories")
	Configs.DB.Table("inventories as i").Select(`i.id,i.name,i.nup,i.year,i.deleted_at,i.updated_at, i.created_at, i.procurement_doc_url, i.status_doc_url, i.image_url,
	i.quantity,i.price,un.name as unit_name, gt.id as goods_type_id, gt.name as goods_type_name, gt.code as goods_type_code, hr.room_id, hc.condition_id, r.name as room_name, c.name as condition_name, us.name as updater_name
	`).Joins("left join (?) as hr on hr.inventory_id = i.id ", subQueryRoom).
		Joins("left join (?) as hc on hc.inventory_id = i.id ", subQueryCondition).
		Joins("left join rooms as r on r.id = hr.room_id").
		Joins("left join conditions as c on c.id = hc.condition_id").
		Joins("left join units as un on un.id = i.unit_id").
		Joins("left join goods_types as gt on gt.id = i.goods_type_id").
		Joins("left join users as us on us.id = i.updater_id").
		Where("i.deleted_at is NULL").Offset(pageSize * (page - 1)).Limit(pageSize).Scan(&inventories)

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

func InventoryPeriod(c *gin.Context) {
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
