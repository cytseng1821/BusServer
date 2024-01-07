package v1

import (
	"BusServer/constant"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Follow stop of citybus by city and route name
// @Produce application/json
// @Param route query string true "route_id"
// @Param route query string true "stop_id"
// @Success 200 {object} constant.Response
// @Router /api/v1/citybus/stop [post]
func FollowCityBusStop(c *gin.Context) {
	var params struct {
		UserID  string `form:"user_id" binding:"required"` // parse from header token (handle at middleware) and add to param
		RouteID string `form:"route_id" binding:"required"`
		StopID  string `form:"stop_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.INVALID_PARAMS, err.Error())
		return
	}

	// [TODO] Upsert subscription relationship into database
	fmt.Println(params)

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}
