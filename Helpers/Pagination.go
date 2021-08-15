package Helpers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type paginateRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var paginData paginateRequest
		var page, pageSize = 1, 10
		if c.Bind(&paginData) == nil {
			if paginData.Page > 0 {
				page = paginData.Page
			}

			if paginData.PageSize > 0 && paginData.PageSize <= 50 {
				pageSize = paginData.PageSize
			}
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
