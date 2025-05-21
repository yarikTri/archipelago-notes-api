package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yarikTri/archipelago-notes-api/internal/common/http/middleware"
	authDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/auth/delivery/http"
)

func InitRoutes(
	authHandler *authDelivery.Handler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	auth := r.Group("/auth-service")
	auth.GET("/login", authHandler.CheckSession)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/registration", authHandler.SignUp)
	auth.POST("/clear-sessions", authHandler.ClearAllSessions)

	return r
}
