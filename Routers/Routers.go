package Routers

import (
	"github.com/DeniesKresna/beinventaris/Controllers"
	"github.com/DeniesKresna/beinventaris/Middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"}, /*
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://github.com"
			},
			MaxAge: 12 * time.Hour,*/
	}))
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/", Middlewares.Auth("administrator"))

		auth.GET("/users", Controllers.UserIndex)
		auth.GET("/users/me", Controllers.UserMe)
		v1.POST("/users", Controllers.UserStore)
		auth.PUT("/users/:id", Controllers.UserUpdate)

		auth.GET("/roles", Controllers.RoleIndex)
		auth.POST("/roles", Controllers.RoleStore)
		auth.PUT("/roles/:id", Controllers.RoleUpdate)

		v1.POST("users/login", Controllers.AuthLogin)

		v1.GET("/academies/list", Controllers.AcademyList)
		v1.GET("/academies/id/:id", Controllers.AcademyShow)
		v1.GET("/academies", Controllers.AcademyIndex)
		v1.POST("/academies", Controllers.AcademyStore)
		v1.PATCH("/academies/:id", Controllers.AcademyUpdate)
		v1.DELETE("/academies/:id", Controllers.AcademyDestroy)

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

		auth.GET("/inventories/list", Controllers.InventoryList)
		auth.GET("/inventories/detail", Controllers.InventoryShow)
		auth.GET("/inventories", Controllers.InventoryIndex)
		auth.POST("/inventories/:id", Controllers.InventoryUpdate)
		auth.POST("/inventories", Controllers.InventoryStore)
		auth.DELETE("/inventories/:id", Controllers.InventoryDestroy)

		v1.GET("/medias", func(c *gin.Context) {
			mediaFile := c.Query("path")
			c.File(mediaFile)
		})

		//v1.GET("users", Controllers.UserIndex)
		//v1.GET("users/:id", Controllers.ShowUser)
		//v1.PUT("users/:id", Controllers.UserUpdate)
		//v1.DELETE("users/:id", Controllers.DestroyUser)
	}
	return r
}
