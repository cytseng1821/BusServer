package v1

import (
	"BusServer/constant"
	"BusServer/controllers/tdx"
	"BusServer/postgresql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Search routes of citybus by city and route name
// @Produce application/json
// @Param city query string true "city"
// @Param route query string true "route"
// @Success 200 {object} constant.Response
// @Router /api/v1/citybus/routes [get]
func SearchCityBusRoutes(c *gin.Context) {
	var params struct {
		City  string `form:"city" binding:"required"`
		Route string `form:"route" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.INVALID_PARAMS, err.Error())
		return
	}

	// [TODO] cache TDX response

	// Get routes info from TDX
	cityBusRoutesInfo, statusCode, err := tdx.GetCityBusRoutes(c, params.City, params.Route)
	if err != nil {
		if statusCode != http.StatusUnauthorized {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
		// refresh token
		if _, err := tdx.GetTDXToken(c); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
		// call again
		if cityBusRoutesInfo, statusCode, err = tdx.GetCityBusRoutes(c, params.City, params.Route); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
	}

	// Retrieve data we need
	var routes []postgresql.CityBusRoute
	for _, route := range cityBusRoutesInfo {
		fmt.Printf("\t%s %s [%s] %s -> %s\n", route.City, route.RouteUID, route.RouteName["Zh_tw"], route.DepartureStopName, route.DestinationStopName)
		for _, subRoute := range route.SubRoutes {
			fmt.Printf("\t\t%s [%s] %v\n", subRoute.SubRouteUID, subRoute.SubRouteName["Zh_tw"], subRoute.Direction)
			routes = append(routes, postgresql.CityBusRoute{
				RouteID:             route.RouteUID,
				RouteName:           route.RouteName["Zh_tw"],
				SubRouteID:          subRoute.SubRouteUID,
				SubRouteName:        subRoute.SubRouteName["Zh_tw"],
				Direction:           subRoute.Direction,
				City:                route.City,
				DepartureStopName:   route.DepartureStopName,
				DestinationStopName: route.DestinationStopName,
			})
		}
	}

	// Upsert routes info into database
	if err := postgresql.InsertCityBusRoutes(c, routes); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_DATABASE, err.Error())
		return
	}

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, routes)
}

// @Summary Search stops of citybus by city and route name
// @Produce application/json
// @Param city query string true "city"
// @Param route query string true "route"
// @Success 200 {object} constant.Response
// @Router /api/v1/citybus/stops [get]
func SearchCityBusStops(c *gin.Context) {
	var params struct {
		City  string `form:"city" binding:"required"`
		Route string `form:"route" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.INVALID_PARAMS, err.Error())
		return
	}

	// [TODO] cache TDX response

	// Get routes stops from TDX
	cityBusRoutes, statusCode, err := tdx.GetCityBusRoutesStops(c, params.City, params.Route)
	if err != nil {
		if statusCode != http.StatusUnauthorized {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
		// refresh token
		if _, err := tdx.GetTDXToken(c); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
		// call again
		if cityBusRoutes, statusCode, err = tdx.GetCityBusRoutesStops(c, params.City, params.Route); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
	}

	// Retrieve data we need
	var stops []postgresql.CityBusStop
	for _, route := range cityBusRoutes {
		fmt.Printf("\t%s %s [%s]: %s [%s] %v\n", route.City, route.RouteUID, route.RouteName["Zh_tw"], route.SubRouteUID, route.SubRouteName["Zh_tw"], route.Direction)
		for _, stop := range route.Stops {
			fmt.Printf("\t\t%v: %s [%s]\n", stop.StopSequence, stop.StopUID, stop.StopName["Zh_tw"])
			stops = append(stops, postgresql.CityBusStop{
				RouteID:      route.RouteUID,
				SubRouteID:   route.SubRouteUID,
				Direction:    route.Direction,
				StopID:       stop.StopUID,
				StopName:     stop.StopName["Zh_tw"],
				StopSequence: stop.StopSequence,
			})
		}
	}

	// Upsert routes stops into database
	if err := postgresql.InsertCityBusRoutesStops(c, stops); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_DATABASE, err.Error())
		return
	}

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, stops)
}
