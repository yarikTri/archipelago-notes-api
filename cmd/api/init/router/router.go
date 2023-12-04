package router

import (
	"github.com/gin-gonic/gin"

	routeDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	ticketDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
)

const TEMPLATE_ROUTE = "static/template/*"

func InitRoutes(
	routeH *routeDelivery.Handler,
	ticketH *ticketDelivery.Handler,
) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob(TEMPLATE_ROUTE)

	r.Static("/static", "./static")
	r.Static("/css", "./static")
	r.Static("/image", "./static")

	r.GET("/routes/:id", routeH.GetByID)
	r.GET("/routes", routeH.List)
	r.GET("/routes/search", routeH.Search)
	r.POST("/routes", routeH.Create)
	r.PUT("/routes/:id", routeH.Update)
	r.DELETE("/routes/:id", routeH.DeleteByID)
	r.PUT("/routes/:id/image", routeH.PutImage)

	r.GET("/tickets/:id", ticketH.GetByID)
	r.GET("/tickets", ticketH.List)
	// r.POST("/tickets", ticketH.Create)
	// r.PUT("/tickets/:id", ticketH.Update)
	r.PUT("/tickets/:id/form", ticketH.FormByID)
	r.PUT("/tickets/:id/moderate", ticketH.ModerateByID)
	r.DELETE("/tickets/:id", ticketH.DeleteByID)

	r.POST("/tickets/routes/:route_id", ticketH.AddRoute)
	r.DELETE("/tickets/:id/routes/:route_id", ticketH.DeleteRoute)

	return r
}
