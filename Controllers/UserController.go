package Controllers

import (
	"strconv"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Helpers"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func UserIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")
	var users []Models.User
	var count int64

	Configs.DB.Model(Models.User{}).Scopes(FilterModel(search, Models.User{})).Count(&count)
	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.
		Preload("Role").Where("id > ?", 0).Scopes(FilterModel(search, Models.User{})), &users, count)
	Response.Json(c, 200, p)
}

func UserStore(c *gin.Context) {
	var user Models.User
	var userCreate Models.UserCreate
	c.ShouldBind(&userCreate)

	v := validate.Struct(userCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	err := Configs.DB.Where("username = ?", userCreate.Username).Or("email = ?", userCreate.Email).First(&Models.User{}).Error
	if err == nil {
		Response.Json(c, 404, "Sudah ada user tersebut")
		return
	}

	if userCreate.Password == "" {
		userCreate.Password = "bawasluindonesia"
	}

	hashedPassword, err := Helpers.Hash(userCreate.Password)
	if err != nil {
		Response.Json(c, 400, "error hashing password")
		return
	}
	InjectStruct(&userCreate, &user)
	user.Password = string(hashedPassword)

	var notAdminRole Models.Role
	err = Configs.DB.Where("name != ?", "administrator").First(&notAdminRole).Error
	if err != nil {
		Response.Json(c, 400, "Harus ada role selain administrator")
		return
	}
	if user.RoleID == 0 {
		user.RoleID = notAdminRole.ID
	}

	if err := Configs.DB.Model(Models.User{}).Create(&user).Error; err != nil {
		Response.Json(c, 500, "Tidak bisa buat pengguna")
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

	id, _ := strconv.Atoi(c.Param("id"))
	var user Models.User

	err := Configs.DB.First(&user, id).Error
	if err != nil {
		Response.Json(c, 404, "Pengguna tidak ditemukan")
		return
	}

	var userUpdate Models.User
	InjectStruct(userUpdateInput, &userUpdate)
	if err := Configs.DB.Model(&user).Updates(userUpdate).Error; err != nil {
		Response.Json(c, 500, "Tidak bisa ubah pengguna")
	} else {
		Response.Json(c, 200, "Berhasil")
	}
}

func UserReset(c *gin.Context) {
	hashedPassword, err := Helpers.Hash("bawasluindonesia")
	if err != nil {
		Response.Json(c, 400, "error hashing password")
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var user Models.User

	err = Configs.DB.First(&user, id).Error
	if err != nil {
		Response.Json(c, 404, "Pengguna tidak ditemukan")
		return
	}

	if err := Configs.DB.Model(&user).Updates(Models.User{Password: string(hashedPassword)}).Error; err != nil {
		Response.Json(c, 500, "Gagal reset password pengguna")
	} else {
		Response.Json(c, 200, "Berhasil")
	}
}

func UserChangePassword(c *gin.Context) {
	SetSessionId(c)

	var userChangePassword Models.UserChangePassword
	var user Models.User
	c.ShouldBindJSON(&userChangePassword)

	v := validate.Struct(userChangePassword)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}

	if err := Configs.DB.First(&user, SessionId).Error; err != nil {
		Response.Json(c, 404, "Pengguna tidak ditemukan")
		return
	}

	err := Helpers.VerifyPassword(user.Password, userChangePassword.OldPassword)
	if err != nil {
		Response.Json(c, 404, "Password salah")
		return
	}

	hashedPassword, err := Helpers.Hash(userChangePassword.Password)
	if err != nil {
		Response.Json(c, 400, "error hashing password")
		return
	}

	if err := Configs.DB.First(&Models.User{}, SessionId).Updates(Models.User{Password: string(hashedPassword)}).Error; err != nil {
		Response.Json(c, 500, "Gagal ganti password pengguna")
	} else {
		Response.Json(c, 200, "Berhasil")
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
