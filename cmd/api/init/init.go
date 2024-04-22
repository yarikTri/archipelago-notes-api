package init

import (
	"github.com/jmoiron/sqlx"
	"github.com/yarikTri/archipelago-notes-api/cmd/api/clients/email"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/archipelago-notes-api/cmd/api/init/router"

	dirsHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/delivery/http"
	dirsRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/repository/postgresql"
	dirsUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/usecase"

	notesHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/delivery/http"
	notesRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/repository/postgresql"
	notesUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/usecase"

	usersHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/delivery/http"
	usersRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/repository/postgresql"
	usersUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/usecase"
)

func Init(sqlDBClient *sqlx.DB, logger logger.Logger) (http.Handler, error) {
	emailInvitationClient := email.NewSmtpInvitationClient()

	notesRepo := notesRepository.NewPostgreSQL(sqlDBClient)
	dirsRepo := dirsRepository.NewPostgreSQL(sqlDBClient)
	usersRepo := usersRepository.NewPostgreSQL(sqlDBClient)

	notesUsecase := notesUsecase.NewUsecase(notesRepo, usersRepo, emailInvitationClient)
	dirsUsecase := dirsUsecase.NewUsecase(dirsRepo, notesRepo)
	usersUsecase := usersUsecase.NewUsecase(usersRepo)

	notesHandler := notesHandler.NewHandler(notesUsecase, logger)
	dirsHandler := dirsHandler.NewHandler(dirsUsecase, logger)
	usersHandler := usersHandler.NewHandler(usersUsecase, logger)

	return router.InitRoutes(
		notesHandler,
		dirsHandler,
		usersHandler,
	), nil
}
