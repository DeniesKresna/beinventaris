package Controllers

import (
	"strconv"
	"strings"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Helpers"
	"github.com/DeniesKresna/beinventaris/Models"
	"github.com/DeniesKresna/beinventaris/Response"
	"github.com/DeniesKresna/beinventaris/Translations"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func AcademyIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var academies []Models.Academy
	p, _ := (&PConfig{
		Page:    page,
		PerPage: pageSize,
		Path:    c.FullPath(),
		Sort:    "id desc",
	}).Paginate(Configs.DB.Preload("Creator"), &academies)
	Response.Json(c, 200, p)
}

func AcademyList(c *gin.Context) {
	var academies []Models.Academy

	Configs.DB.Find(&academies)
	Response.Json(c, 200, academies)
}

func AcademyShow(c *gin.Context) {
	id := c.Param("id")
	var academy Models.Academy
	err := Configs.DB.Preload("Creator").First(&academy, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.AcademyNotFound)
		return
	}

	Response.Json(c, 200, academy)
}

func AcademyStore(c *gin.Context) {
	SetSessionId(c)
	var academy Models.Academy
	var academyCreate Models.AcademyCreate

	//bind and validate request-------------------------
	if err := c.ShouldBind(&academyCreate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(academyCreate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.Where("name = ?", academyCreate.Name).First(&Models.Academy{}).Error
	if err == nil {
		Response.Json(c, 500, Translations.AcademyExist)
		return
	}

	academyCreate.CreatorID = SessionId

	InjectStruct(&academyCreate, &academy)
	if err := Configs.DB.Create(&academy).Error; err != nil {
		Response.Json(c, 500, Translations.AcademyCreateServerError)
		return
	} else {
		file, err := c.FormFile("image")
		if err != nil {
			Configs.DB.Unscoped().Delete(&academy)
			Response.Json(c, 500, Translations.AcademyCreateUploadError)
			return
		}
		filename := "academy-" + strconv.FormatUint(uint64(academy.ID), 10) + "-" + file.Filename
		filename = strings.ReplaceAll(filename, " ", "-")
		if err := c.SaveUploadedFile(file, Helpers.AcademyPath(filename)); err != nil {
			Configs.DB.Unscoped().Delete(&academy)
			Response.Json(c, 500, Translations.AcademyCreateUploadError)
			return
		}
		if err := Configs.DB.Model(&academy).Update("image_url", Helpers.AcademyPath(filename)).Error; err != nil {
			Configs.DB.Unscoped().Delete(&academy)
			Response.Json(c, 500, Translations.AcademyCreateUploadError)
			return
		}

		Response.Json(c, 200, Translations.AcademyCreated)
	}
}

func AcademyUpdate(c *gin.Context) {
	SetSessionId(c)
	var academy Models.Academy
	var academyUpdate Models.AcademyUpdate
	id := c.Param("id")

	//bind and validate request-------------------------
	if err := c.ShouldBind(&academyUpdate); err != nil {
		Response.Json(c, 422, err)
		return
	}
	v := validate.Struct(academyUpdate)
	if !v.Validate() {
		Response.Json(c, 422, v.Errors.One())
		return
	}
	//--------------------------------------------------

	err := Configs.DB.First(&academy, id).Error
	if err != nil {
		Response.Json(c, 404, Translations.AcademyNotFound)
		return
	}

	academyUpdate.CreatorID = SessionId
	InjectStruct(&academyUpdate, &academy)
	if err := Configs.DB.Save(&academy).Error; err != nil {
		Response.Json(c, 500, Translations.AcademyUpdateServerError)
		return
	} else {
		file, err := c.FormFile("image")
		if err == nil {
			filename := "academy-" + strconv.FormatUint(uint64(academy.ID), 10) + "-" + file.Filename
			filename = strings.ReplaceAll(filename, " ", "-")
			if err := c.SaveUploadedFile(file, Helpers.AcademyPath(filename)); err != nil {
				Configs.DB.Unscoped().Delete(&academy)
				Response.Json(c, 500, Translations.AcademyUpdateUploadError)
				return
			}
			oldfilename := academy.ImageUrl
			if err := Configs.DB.Model(&academy).Update("image_url", Helpers.AcademyPath(filename)).Error; err != nil {
				Response.Json(c, 500, Translations.AcademyUpdateUploadError)
				return
			}
			err := Helpers.DeleteFile(oldfilename)
			if err != nil {
			}
		} else {
			Response.Json(c, 500, Translations.AcademyUpdateUploadError)
			return
		}

		Response.Json(c, 200, Translations.AcademyUpdated)
	}
}

func AcademyDestroy(c *gin.Context) {
	id := c.Param("id")
	var academy Models.Academy
	err := Configs.DB.First(&academy, id).Error

	if err != nil {
		Response.Json(c, 404, Translations.AcademyNotFound)
		return
	}

	Configs.DB.Delete(&academy)

	Response.Json(c, 404, Translations.AcademyDeleted)
	return
}
