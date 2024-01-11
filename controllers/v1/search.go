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
	var subRoutes []postgresql.CityBusSubRoute
	for _, r := range cityBusRoutesInfo {
		fmt.Printf("\t%s %s [%s] %s -> %s\n", r.City, r.RouteUID, r.RouteName["Zh_tw"], r.DepartureStopName, r.DestinationStopName)
		routes = append(routes, postgresql.CityBusRoute{
			RouteID:             r.RouteUID,
			RouteName:           r.RouteName["Zh_tw"],
			City:                r.City,
			DepartureStopName:   r.DepartureStopName,
			DestinationStopName: r.DestinationStopName,
		})
		for _, sr := range r.SubRoutes {
			fmt.Printf("\t\t%s [%s] %v\n", sr.SubRouteUID, sr.SubRouteName["Zh_tw"], sr.Direction)
			subRoutes = append(subRoutes, postgresql.CityBusSubRoute{
				RouteID:      r.RouteUID,
				SubRouteID:   sr.SubRouteUID,
				SubRouteName: sr.SubRouteName["Zh_tw"],
				Direction:    sr.Direction,
			})
		}
	}

	// Upsert routes info into database
	if err := postgresql.InsertCityBusRoutes(c, routes, subRoutes); err != nil {
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
		City      string `form:"city" binding:"required"`  // for TDX api
		Route     string `form:"route" binding:"required"` // for TDX api
		RouteID   string `form:"route_id" binding:"required"`
		Direction *int   `form:"direction" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.INVALID_PARAMS, err.Error())
		return
	}

	// [TODO] cache TDX response

	// Get routes stops from TDX
	cityBusSubRoutes, statusCode, err := tdx.GetCityBusRoutesStops(c, params.City, params.Route)
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
		if cityBusSubRoutes, statusCode, err = tdx.GetCityBusRoutesStops(c, params.City, params.Route); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_TDX, err.Error())
			return
		}
	}

	// Retrieve data we need
	stopMap := map[string]bool{}
	var stops []postgresql.CityBusStop
	var relations []postgresql.CityBusSubRoute2StopRelation
	for _, sr := range cityBusSubRoutes {
		if sr.RouteUID != params.RouteID || sr.Direction != *params.Direction {
			continue
		}
		fmt.Printf("\t%s %s [%s]: %s [%s] %v\n", sr.City, sr.RouteUID, sr.RouteName["Zh_tw"], sr.SubRouteUID, sr.SubRouteName["Zh_tw"], sr.Direction)
		for _, s := range sr.Stops {
			fmt.Printf("\t\t%v: %s [%s]\n", s.StopSequence, s.StopUID, s.StopName["Zh_tw"])
			if !stopMap[s.StopUID] { // de-duplicate before input due to 'ON CONFLICT' operation
				stops = append(stops, postgresql.CityBusStop{
					RouteID:  sr.RouteUID,
					StopID:   s.StopUID,
					StopName: s.StopName["Zh_tw"],
				})
				stopMap[s.StopUID] = true
			}
			relations = append(relations, postgresql.CityBusSubRoute2StopRelation{
				SubRouteID:   sr.SubRouteUID,
				Direction:    sr.Direction,
				StopID:       s.StopUID,
				StopSequence: s.StopSequence,
			})
		}
	}

	// Upsert routes stops into database
	if err := postgresql.InsertCityBusRoutesStops(c, params.RouteID, stops, relations); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_DATABASE, err.Error())
		return
	}

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, stops)
}
