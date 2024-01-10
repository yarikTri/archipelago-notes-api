package init

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/redis"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/s3/minio"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/router"

	routeHandler "github.com/yarikTri/web-transport-cards/internal/pkg/route/delivery/http"
	routeRepository "github.com/yarikTri/web-transport-cards/internal/pkg/route/repository/postgresql"
	routeUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/route/usecase"

	imageRepository "github.com/yarikTri/web-transport-cards/internal/pkg/image/repository/s3"
	imageUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/image/usecase"

	userRepository "github.com/yarikTri/web-transport-cards/internal/pkg/user/repository/postgresql"

	authHandler "github.com/yarikTri/web-transport-cards/internal/pkg/auth/delivery/http"
	authRepository "github.com/yarikTri/web-transport-cards/internal/pkg/auth/repository/redis"
	authUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/auth/usecase"

	ticketHandler "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/delivery/http"
	ticketRepository "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/repository/postgresql"
	ticketDraftRepository "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/repository/redis"
	ticketUsecase "github.com/yarikTri/web-transport-cards/internal/pkg/ticket/usecase"
)

const (
	MINIO_DIR = "images"
)

func Init(sqlDBClient *gorm.DB, logger logger.Logger) (http.Handler, error) {
	s3MinioClient, err := minio.MakeS3MinioClient()
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.InitRedisDB()
	if err != nil {
		return nil, err
	}

	userRepo := userRepository.NewPostgreSQL(sqlDBClient)

	routeRepo := routeRepository.NewPostgreSQL(sqlDBClient)
	ticketRepo := ticketRepository.NewPostgreSQL(sqlDBClient)
	ticketDraftRepo := ticketDraftRepository.NewRedis(redisClient)
	imageRepo := imageRepository.NewS3MinioImageStorage(MINIO_DIR, s3MinioClient)
	authRepo := authRepository.NewRedis(redisClient)

	routeUsecase := routeUsecase.NewUsecase(routeRepo)
	ticketUsecase := ticketUsecase.NewUsecase(ticketRepo, ticketDraftRepo, routeRepo)
	imageUsecase := imageUsecase.NewUsecase(imageRepo)
	authUsecase := authUsecase.NewUsecase(authRepo, userRepo)

	routeHandler := routeHandler.NewHandler(routeUsecase, imageUsecase, ticketUsecase, authUsecase, logger)
	ticketHandler := ticketHandler.NewHandler(ticketUsecase, authUsecase, logger)
	authHandler := authHandler.NewHandler(authUsecase, logger)

	return router.InitRoutes(
		routeHandler,
		ticketHandler,
		authHandler,
	), nil
}
