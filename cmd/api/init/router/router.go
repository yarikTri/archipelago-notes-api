package router

import (
	"github.com/gin-gonic/gin"

	routeDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	stationDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/station/delivery/http"
	ticketDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
)

const TEMPLATE_ROUTE = "static/template/*"

func InitRoutes(
	routeH *routeDelivery.Handler,
	stationH *stationDelivery.Handler,
	ticketH *ticketDelivery.Handler,
) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob(TEMPLATE_ROUTE)

	r.Static("/static", "./static")
	r.Static("/css", "./static")
	r.Static("/image", "./static")

	r.GET("/", routeH.List)

	r.GET("/search", routeH.Search)
	r.GET("/routes/:id", routeH.GetByID)

	return r
}
