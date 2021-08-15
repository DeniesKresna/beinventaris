package Routers

import (
	"github.com/DeniesKresna/jobhunop/Controllers"
	"github.com/DeniesKresna/jobhunop/Middlewares"
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
		auth.POST("/academies", Controllers.AcademyStore)
		auth.PATCH("/academies/:id", Controllers.AcademyUpdate)
		auth.DELETE("/academies/:id", Controllers.AcademyDestroy)

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
