package init

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/yarikTri/archipelago-notes-api/internal/clients/invitations/email"
	"github.com/yarikTri/archipelago-notes-api/internal/clients/llm"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/archipelago-notes-api/cmd/api/init/router"

	dirsHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/delivery/http"
	dirsRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/repository/postgresql"
	dirsUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs/usecase"

	notesHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/delivery/http"
	notesRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/repository/postgresql"
	notesUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/notes/usecase"

	tagHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/delivery/http"
	tagSuggester "github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/repository/ollama"
	tagRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/repository/postgresql"
	tagUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/usecase"

	usersHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/delivery/http"
	usersRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/repository/postgresql"
	usersUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/users/usecase"

	summaryHandler "github.com/yarikTri/archipelago-notes-api/internal/pkg/summary/delivery/http"
	summaryRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/summary/repository/postgresql"
	summaryUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/summary/usecase"
)

func Init(sqlDBClient *sqlx.DB, logger logger.Logger, openAiUrl string, tagSuggesterModel string, defaultGenerateTagNum int) (http.Handler, error) {
	emailClient := email.NewEmailClient()
	openAiClient := llm.NewOpenAiClient(openAiUrl)

	tagSuggesterRepo := tagSuggester.NewTagSuggester(openAiClient, defaultGenerateTagNum, tagSuggesterModel)

	notesRepo := notesRepository.NewPostgreSQL(sqlDBClient)
	dirsRepo := dirsRepository.NewPostgreSQL(sqlDBClient)
	usersRepo := usersRepository.NewPostgreSQL(sqlDBClient)
	summRepo := summaryRepository.NewPostgreSQL(sqlDBClient)
	tagRepo := tagRepository.NewPostgreSQL(sqlDBClient)

	notesUsecase := notesUsecase.NewUsecase(notesRepo, usersRepo, emailClient)
	dirsUsecase := dirsUsecase.NewUsecase(dirsRepo, notesRepo)
	usersUsecase := usersUsecase.NewUsecase(usersRepo, emailClient)
	summaryUsecase := summaryUsecase.NewUsecase(summRepo)
	tagUsecase := tagUsecase.NewUsecase(tagRepo, tagSuggesterRepo)

	notesHandler := notesHandler.NewHandler(notesUsecase, logger)
	dirsHandler := dirsHandler.NewHandler(dirsUsecase, logger)
	usersHandler := usersHandler.NewHandler(usersUsecase, logger)
	summaryHandler := summaryHandler.NewHandler(summaryUsecase, logger)
	tagHandler := tagHandler.NewHandler(tagUsecase, notesUsecase, logger)

	return router.InitRoutes(
		notesHandler,
		dirsHandler,
		usersHandler,
		summaryHandler,
		tagHandler,
	), nil
}
