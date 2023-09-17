package init

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/postgresql"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/router"

	routeHandler "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	routeRepository "github.com/yarikTri/web-transport-cards/internal/pkg/route/repository/postgresql"
	routeUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/route/usecase"
	stationHandler "github.com/yarikTri/web-transport-cards/internal/pkg/station/delivery/http"
	stationRepository "github.com/yarikTri/web-transport-cards/internal/pkg/station/repository/postgresql"
	stationUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/station/usecase"
	ticketHandler "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
	ticketRepository "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/repository/postgresql"
	ticketUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/usecase"
)

func Init(db *gorm.DB, tables postgresql.PostgreSQLTables, logger logger.Logger) (http.Handler, error) {
	routeRepo := routeRepository.NewPostgreSQL(db, tables)
	stationRepo := stationRepository.NewPostgreSQL(db, tables)
	ticketRepo := ticketRepository.NewPostgreSQL(db, tables)

	routeUsecase := routeUsecase.NewUsecase(routeRepo)
	stationUsecase := stationUsecase.NewUsecase(stationRepo)
	ticketUsecase := ticketUsecase.NewUsecase(ticketRepo)

	routeHandler := routeHandler.NewHandler(routeUsecase, logger)
	stationHandler := stationHandler.NewHandler(stationUsecase, logger)
	ticketHandler := ticketHandler.NewHandler(ticketUsecase, logger)

	return router.InitRoutes(
		routeHandler,
		stationHandler,
		ticketHandler,
	), nil
}
