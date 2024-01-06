package v1

import (
	"BusServer/constant"
	"BusServer/controllers/tdx"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Refresh access token for TDX api
// @Produce application/json
// @Success 200 {object} constant.Response
// @Router /api/v1/token [get]
func RefreshTDXToken(c *gin.Context) {
	accessToken, err := tdx.GetTDXToken(c)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
		return
	}

	fmt.Println(accessToken)
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}
