package Middlewares

import (
	"os"
	"strings"

	"github.com/DeniesKresna/jobhunop/Configs"
	"github.com/DeniesKresna/jobhunop/Models"
	"github.com/DeniesKresna/jobhunop/Response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, BEARER_SCHEMA) {
			Response.Json(c, 401, "You must login first")
			c.Abort()
			return
		}
		tokenString := authHeader[(len(BEARER_SCHEMA) + 1):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			apiSecret := os.Getenv("API_SECRET")
			return []byte(apiSecret), nil
		})

		if err != nil {
			Response.Json(c, 401, "Not allowed access this page")
			c.Abort()
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			userId := claims["userId"]
			userId = uint(userId.(float64))

			var user Models.User
			err = Configs.DB.First(&user, userId).Joins("INNER JOIN roles r ON user.role_id = r.id").
				Where("r.name = ?", "admin").Error
			if err != nil {
				Response.Json(c, 404, "User not found")
				c.Abort()
				return
			}

			c.Set("sessionId", userId)
			//fmt.Println(claims)
		} else {
			Response.Json(c, 401, "You dont have access")
			c.Abort()
			return
		}
	}
}
