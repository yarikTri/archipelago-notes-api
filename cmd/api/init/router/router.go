package router

import (
	"github.com/gin-gonic/gin"

	routeDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	stationDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/station/delivery/http"
	ticketDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
)

func InitRoutes(
	routeH *routeDelivery.Handler,
	stationH *stationDelivery.Handler,
	ticketH *ticketDelivery.Handler) *gin.Engine {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", ticketH.Index)

	r.GET("/routes/get", routeH.GetByID)
	r.GET("/routes/list", routeH.List)
	r.GET("/routes/create", routeH.Create)
	r.GET("/routes/delete", routeH.DeleteByID)
	r.GET("/stations/get", stationH.GetByID)
	r.GET("/stations/list", stationH.List)
	r.GET("/stations/create", stationH.Create)
	r.GET("/stations/delete", stationH.DeleteByID)
	r.GET("/tickets/get", ticketH.GetByID)
	r.GET("/tickets/list", ticketH.List)
	r.GET("/tickets/create", ticketH.Create)
	r.GET("/tickets/delete", ticketH.DeleteByID)

	r.Static("/images", "./resources/images")
	r.Static("/docs", "./resources/docs")

	return r
}
