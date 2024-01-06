package routes

import (
	"BusServer/controllers"
	v1 "BusServer/controllers/v1"
	"BusServer/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowWildcard = true
	config.AddAllowMethods(http.MethodOptions)
	config.AddAllowHeaders("Authorization", "Content-Type", "Upgrade", "Connection", "Accept", "Accept-Encoding", "Accept-Language", "Host", "Cookie", "Referer", "User-Agent")
	if err := config.Validate(); err != nil {
		panic(err)
	}
	router.Use(cors.New(config))

	router.GET("/heartBeat", controllers.HeartBeat) // check alive

	apiv1 := router.Group("/api/v1")
	apiv1.Use(middleware.MiddleWare)

	apiv1.GET("/citybus/routes", v1.SearchCityBusRoutes)
	apiv1.GET("/citybus/stops", v1.SearchCityBusStops)
	apiv1.POST("/citybus/stop", v1.FollowCityBusStop)

	apiv1.POST("/token", v1.RefreshTDXToken)

	return router
}
