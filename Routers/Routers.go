package Routers

import (
	"github.com/DeniesKresna/beinventaris/Controllers"
	"github.com/DeniesKresna/beinventaris/Middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Disposition"}, /*
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://github.com"
			},
			MaxAge: 12 * time.Hour,*/
	}))
	// Serve frontend static files
	r.Use(static.Serve("/", static.LocalFile("./Client/Build", true)))
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/", Middlewares.Auth("administrator"))

		auth.GET("/users", Controllers.UserIndex)
		auth.GET("/users/me", Controllers.UserMe)
		auth.GET("/users/reset/:id", Controllers.UserReset)
		auth.POST("/users", Controllers.UserStore)
		auth.POST("/users/change-password", Controllers.UserChangePassword)
		auth.PATCH("/users/:id", Controllers.UserUpdate)

		auth.GET("/roles", Controllers.RoleIndex)
		auth.GET("/roles/list", Controllers.RoleList)
		auth.POST("/roles", Controllers.RoleStore)
		auth.PUT("/roles/:id", Controllers.RoleUpdate)

		v1.POST("users/login", Controllers.AuthLogin)

		auth.GET("/units/list", Controllers.UnitList)
		auth.GET("/units/id/:id", Controllers.UnitShow)
		auth.GET("/units", Controllers.UnitIndex)
		auth.POST("/units", Controllers.UnitStore)
		auth.DELETE("/units/:id", Controllers.UnitDestroy)

		auth.GET("/rooms/list", Controllers.RoomList)
		auth.GET("/rooms/id/:id", Controllers.RoomShow)
		auth.GET("/rooms", Controllers.RoomIndex)
		auth.POST("/rooms", Controllers.RoomStore)
		auth.DELETE("/rooms/:id", Controllers.RoomDestroy)

		auth.GET("/goods-types/list", Controllers.GoodsTypeList)
		auth.GET("/goods-types/id/:id", Controllers.GoodsTypeShow)
		auth.GET("/goods-types", Controllers.GoodsTypeIndex)
		auth.POST("/goods-types", Controllers.GoodsTypeStore)
		auth.PATCH("/goods-types/:id", Controllers.GoodsTypeUpdate)
		auth.DELETE("/goods-types/:id", Controllers.GoodsTypeDestroy)

		auth.GET("/conditions/list", Controllers.ConditionList)
		auth.GET("/conditions/id/:id", Controllers.ConditionShow)
		auth.GET("/conditions", Controllers.ConditionIndex)
		auth.POST("/conditions", Controllers.ConditionStore)
		auth.DELETE("/conditions/:id", Controllers.ConditionDestroy)

		auth.GET("/periods/list", Controllers.PeriodList)
		auth.GET("/periods/id/:id", Controllers.PeriodShow)
		auth.GET("/periods/active", Controllers.PeriodActive)
		auth.GET("/periods", Controllers.PeriodIndex)
		auth.POST("/periods", Controllers.PeriodStore)
		auth.PATCH("/periods/:id", Controllers.PeriodUpdate)
		auth.DELETE("/periods/:id", Controllers.PeriodDestroy)

		auth.GET("/inventories/getexcel", Controllers.InventoryExport)
		auth.GET("/inventories/list", Controllers.InventoryList)
		auth.POST("/inventories-code/detail", Controllers.InventoryShow)
		auth.POST("/inventories-codename/detail", Controllers.InventoryCodeNameShow)
		auth.POST("/inventories-period", Controllers.InventoryUpdatePeriod)
		auth.GET("/inventories/detail/:id", Controllers.InventoryShowDetail)
		auth.GET("/inventories", Controllers.InventoryIndex)
		auth.POST("/inventories/:id", Controllers.InventoryUpdate)
		auth.POST("/inventories", Controllers.InventoryStore)
		auth.DELETE("/inventories/:id", Controllers.InventoryDestroy)
		auth.GET("/inventories-period", Controllers.InventoryPeriodIndex)
		auth.GET("/inventories-period/getexcel", Controllers.InventoryPeriodExport)
		auth.POST("/inventories-period/delete", Controllers.InventoryPeriodDelete)

		auth.GET("/histories/list", Controllers.HistoryList)
		auth.GET("/histories/inventoryid/:id", Controllers.HistoryIndex)
		auth.POST("/histories/:id", Controllers.HistoryUpdate)
		auth.POST("/histories", Controllers.HistoryStore)
		auth.DELETE("/histories/:id", Controllers.HistoryDestroy)

		v1.GET("/medias", func(c *gin.Context) {
			mediaFile := c.Query("path")
			c.File(mediaFile)
		})

		v1.GET("/documents", Controllers.DownloadDocuments)

		//v1.GET("users", Controllers.UserIndex)
		//v1.GET("users/:id", Controllers.ShowUser)
		//v1.PUT("users/:id", Controllers.UserUpdate)
		//v1.DELETE("users/:id", Controllers.DestroyUser)
	}
	return r
}
