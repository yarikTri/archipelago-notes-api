package init

import (
	"github.com/jmoiron/sqlx"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/archipelago-notes-api/cmd/api/init/router"

	notesHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/delivery/http"
	notesRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/repository/postgresql"
	notesUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/usecase"
)

func Init(sqlDBClient *sqlx.DB, logger logger.Logger) (http.Handler, error) {
	notesRepo := notesRepository.NewPostgreSQL(sqlDBClient)

	notesUsecase := notesUsecase.NewUsecase(notesRepo)

	notesHandler := notesHandler.NewHandler(notesUsecase, logger)

	return router.InitRoutes(
		notesHandler,
	), nil
}
