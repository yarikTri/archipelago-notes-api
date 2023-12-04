package init

import (
	"fmt"
	"net/http"
	"os"

	"gorm.io/gorm"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/s3/minio"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/router"

	routeHandler "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	routeRepository "github.com/yarikTri/web-transport-cards/internal/pkg/route/repository/postgresql"
	routeUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/route/usecase"

	imageRepository "github.com/yarikTri/web-transport-cards/internal/pkg/image/repository/s3"
	imageUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/image/usecase"

	ticketHandler "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
	ticketRepository "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/repository/postgresql"
	ticketUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/usecase"
)

const (
	MINIO_DIR = "images"
)

func Init(db *gorm.DB, logger logger.Logger) (http.Handler, error) {
	s3MinioClient, err := minio.MakeS3MinioClient(
		os.Getenv("MINIO_LISTEN_ENDPOINT"), os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"),
	)
	if err != nil {
		return nil, err
	}
	fmt.Println("Minio client is alive: ", s3MinioClient.IsOnline())

	routeRepo := routeRepository.NewPostgreSQL(db)
	imageRepo := imageRepository.NewS3MinioImageStorage(MINIO_DIR, s3MinioClient)
	ticketRepo := ticketRepository.NewPostgreSQL(db)

	routeUsecase := routeUsecase.NewUsecase(routeRepo)
	imageUsecase := imageUsecase.NewUsecase(imageRepo)
	ticketUsecase := ticketUsecase.NewUsecase(ticketRepo)

	routeHandler := routeHandler.NewHandler(routeUsecase, imageUsecase, logger)
	ticketHandler := ticketHandler.NewHandler(ticketUsecase, logger)

	return router.InitRoutes(
		routeHandler,
		ticketHandler,
	), nil
}
