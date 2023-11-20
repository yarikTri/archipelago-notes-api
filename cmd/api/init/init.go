package init

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/router"

	mockDB "github.com/yarikTri/web-transport-cards/cmd/api/init/db/mock"
	routeHandler "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	routeRepository "github.com/yarikTri/web-transport-cards/internal/pkg/route/repository/postgresql"
	routeUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/route/usecase"
	stationHandler "github.com/yarikTri/web-transport-cards/internal/pkg/station/delivery/http"
	mockStationRepository "github.com/yarikTri/web-transport-cards/internal/pkg/station/repository/mock"
	stationUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/station/usecase"
	ticketHandler "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
	ticketRepository "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/repository/postgresql"
	ticketUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/usecase"
)

func Init(db *gorm.DB, logger logger.Logger) (http.Handler, error) {
	routeRepo := routeRepository.NewPostgreSQL(db)
	stationRepo := mockStationRepository.NewMock(mockDB.MockDBImpl)
	ticketRepo := ticketRepository.NewPostgreSQL(db)

	routeUsecase := routeUsecase.NewUsecase(routeRepo)
	stationUsecase := stationUsecase.NewUsecase(stationRepo)
	ticketUsecase := ticketUsecase.NewUsecase(ticketRepo)

	routeHandler := routeHandler.NewHandler(routeUsecase, stationUsecase, logger)
	stationHandler := stationHandler.NewHandler(stationUsecase, logger)
	ticketHandler := ticketHandler.NewHandler(ticketUsecase, logger)

	return router.InitRoutes(
		routeHandler,
		stationHandler,
		ticketHandler,
	), nil
}
