package Controllers

import (
	"strconv"

	"github.com/DeniesKresna/jobhunop/Configs"
	"github.com/DeniesKresna/jobhunop/Models"
	"github.com/DeniesKresna/jobhunop/Response"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func RoleIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var roles []Models.Role
	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB, &roles)
	Response.Json(c, 200, p)
}

func RoleStore(c *gin.Context) {
	SetSessionId(c)
	var role Models.Role
	c.ShouldBindJSON(&role)

	v := validate.Struct(role)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	err := Configs.DB.Where("name = ?", role.Name).First(&Models.Role{}).Error
	if err == nil {
		Response.Json(c, 404, "Sudah ada role tersebut")
		return
	}

	role.CreatorID = SessionId

	if err := Configs.DB.Model(Models.Role{}).Create(&role).Error; err != nil {
		Response.Json(c, 500, "Cant Create Role")
	} else {
		Response.Json(c, 200, "Success")
	}
}

func RoleUpdate(c *gin.Context) {
	SetSessionId(c)
	var roleUpdateInput Models.RoleUpdate
	c.ShouldBindJSON(&roleUpdateInput)
	v := validate.Struct(roleUpdateInput)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	roleUpdateInput.CreatorID = SessionId

	if err := Configs.DB.First(&Models.Role{}, c.Param("id")).Updates(&roleUpdateInput).Error; err != nil {
		Response.Json(c, 500, "Cant Update Role")
	} else {
		Response.Json(c, 200, "Success")
	}
}
