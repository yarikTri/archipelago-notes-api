package router

import (
	"github.com/gin-gonic/gin"

	_ "github.com/yarikTri/web-transport-cards/docs"

	authDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/auth/delivery/http"
	routeDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	ticketDelivery "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"

	middleware "github.com/yarikTri/web-transport-cards/internal/common/http/middleware"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	swagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(
	routeH *routeDelivery.Handler,
	ticketH *ticketDelivery.Handler,
	authH *authDelivery.Handler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	routes := r.Group("/routes")
	routes.GET("/:id", routeH.GetByID)
	routes.GET("", routeH.List)
	routes.POST("", routeH.Create)
	routes.PUT("/:id", routeH.Update)
	routes.DELETE("/:id", routeH.DeleteByID)
	routes.PUT("/:id/image", routeH.PutImage)

	tickets := r.Group("/tickets")
	tickets.GET("/:id", ticketH.GetByID)
	tickets.GET("", ticketH.List)
	tickets.PUT("/:id/moderate", ticketH.ModerateByID)
	tickets.DELETE("/:id", ticketH.DeleteByID)

	// ticket draft
	tickets.GET("/draft", ticketH.GetDraft)
	tickets.DELETE("/draft", ticketH.DeleteDraft)
	tickets.PUT("/draft/form", ticketH.FormDraft)

	// routes M:N tickets
	tickets.POST("/routes/:route_id", ticketH.AddRoute)
	tickets.DELETE("/routes/:route_id", ticketH.DeleteRoute)

	auth := r.Group("/auth")
	auth.POST("/signup", authH.SignUp)
	auth.POST("/login", authH.Login)
	auth.POST("/logout", authH.Logout)

	r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	return r
}
