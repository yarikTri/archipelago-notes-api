package router

import (
	"github.com/gin-gonic/gin"

	_ "github.com/yarikTri/archipelago-notes-api/docs"

	"github.com/yarikTri/archipelago-notes-api/internal/common/http/middleware"
	dirsDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/delivery/http"
	notesDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/delivery/http"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	swagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(
	notesHandler *notesDelivery.Handler,
	dirsHandler *dirsDelivery.Handler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")

	notes := api.Group("/notes")
	notes.GET("/:id", notesHandler.Get)
	notes.GET("", notesHandler.List)
	notes.POST("", notesHandler.Create)
	notes.POST("/:id", notesHandler.Update)
	notes.DELETE("/:id", notesHandler.Delete)

	dirs := api.Group("/dirs")
	dirs.GET("/:id", dirsHandler.Get)
	dirs.GET("/:id/tree", dirsHandler.GetTree)
	dirs.POST("", dirsHandler.Create)
	dirs.POST("/:id", dirsHandler.Update)
	dirs.DELETE("/:id", dirsHandler.Delete)

	r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	return r
}
