package router

import (
	"github.com/gin-gonic/gin"

	_ "github.com/yarikTri/archipelago-notes-api/docs"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	swagger "github.com/swaggo/gin-swagger"
	"github.com/yarikTri/archipelago-notes-api/internal/common/http/middleware"
	dirsDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/delivery/http"
	notesDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/delivery/http"
	summaryDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/summary/delivery/http"
	usersDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/delivery/http"
)

func InitRoutes(
	notesHandler *notesDelivery.Handler,
	dirsHandler *dirsDelivery.Handler,
	usersHandler *usersDelivery.Handler,
	summaryHandler *summaryDelivery.Handler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")

	notes := api.Group("/notes")
	notes.GET("/:id", notesHandler.Get)
	notes.GET("/:id/summaries", notesHandler.GetSummaryListByNote)
	notes.GET("", notesHandler.List)
	notes.POST("", notesHandler.Create)
	notes.POST("/:id", notesHandler.Update)
	notes.DELETE("/:id", notesHandler.Delete)
	notes.POST("/:id/access/:userID", notesHandler.SetAccess)
	notes.GET("/:id/is_owner/:userID", notesHandler.CheckOwner)
	notes.POST("/:id/attach_summ/:summID", notesHandler.AttachNoteToSummary)
	notes.POST("/:id/detach_summ/:summID", notesHandler.DetachNoteFromSummary)
	notes.GET("/:id/summary_list", notesHandler.GetSummaryListByNote)

	dirs := api.Group("/dirs")
	dirs.GET("/:id", dirsHandler.Get)
	dirs.GET("/:id/tree", dirsHandler.GetTree)
	dirs.POST("", dirsHandler.Create)
	dirs.POST("/:id", dirsHandler.Update)
	dirs.DELETE("/:id", dirsHandler.Delete)

	users := api.Group("/users")
	users.GET("/:id", usersHandler.Get)
	users.GET("", usersHandler.Search)
	users.POST("/:userID/root_dir/:rootDirID", usersHandler.SetRootDirID)

	summary := api.Group("/summary")
	summary.GET("/get/:id", summaryHandler.GetSummary)
	summary.GET("/finish/:id", summaryHandler.FinishSummary)
	summary.POST("/save", summaryHandler.SaveSummaryText)
	summary.POST("/update_text_role", summaryHandler.UpdateSummaryTextRole)
	summary.GET("/active", summaryHandler.GetActiveSummaries)

	r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	return r
}
