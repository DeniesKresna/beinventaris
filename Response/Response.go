package Response

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status  int
	Message string
	Data    interface{}
}

func Json(w *gin.Context, status int, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Message = ""
	_, ok := payload.(string)
	if ok {
		res.Message = payload.(string)
	}
	res.Data = payload

	w.JSON(status, res)
}
