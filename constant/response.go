package constant

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS        = http.StatusOK
	INVALID_PARAMS = http.StatusBadRequest
	ERROR          = http.StatusInternalServerError
)

const (
	ERROR_DATABASE = iota + 10000
	ERROR_TDX
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func ResponseWithData(c *gin.Context, httpCode, respCode int, data interface{}) {
	response := Response{
		Code: respCode,
		Data: data,
	}
	c.JSON(httpCode, response)
}
