package Controllers

import (
	"strconv"

	"github.com/DeniesKresna/jobhunop/Configs"
	"github.com/DeniesKresna/jobhunop/Helpers"
	"github.com/DeniesKresna/jobhunop/Models"
	"github.com/DeniesKresna/jobhunop/Response"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func UserIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var users []Models.User
	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.
		Preload("Role").Where("id > ?", 0), &users)
	Response.Json(c, 200, p)
}

func UserStore(c *gin.Context) {
	var user Models.User
	c.ShouldBindJSON(&user)

	v := validate.Struct(user)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	err := Configs.DB.Where("username = ?", user.Username).Or("email = ?", user.Email).First(&Models.User{}).Error
	if err == nil {
		Response.Json(c, 404, "Sudah ada user tersebut")
		return
	}

	hashedPassword, err := Helpers.Hash(user.Password)
	if err != nil {
		Response.Json(c, 400, "error hashing password")
		return
	}
	user.Password = string(hashedPassword)
	user.RoleID = 1

	if err := Configs.DB.Model(Models.User{}).Create(&user).Error; err != nil {
		Response.Json(c, 500, "Cant Create User")
	} else {
		Response.Json(c, 200, "Success")
	}
}

func UserUpdate(c *gin.Context) {
	var userUpdateInput Models.UserUpdate
	c.ShouldBindJSON(&userUpdateInput)
	v := validate.Struct(userUpdateInput)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	if err := Configs.DB.First(&Models.User{}, c.Param("id")).Updates(&userUpdateInput).Error; err != nil {
		Response.Json(c, 500, "Cant Update User")
	} else {
		Response.Json(c, 200, "Success")
	}
}

func UserMe(c *gin.Context) {
	Response.Json(c, 200, me(c))
}

func me(c *gin.Context) *Models.User {
	SetSessionId(c)
	var user Models.User
	Configs.DB.Preload("Role").First(&user, SessionId)
	return &user
}
