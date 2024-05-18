package init

import (
	"fmt"
	"github.com/yarikTri/archipelago-notes-api/cmd/auth/init/db/redis"
	"github.com/yarikTri/archipelago-notes-api/cmd/common/init/db/postgresql"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/archipelago-notes-api/cmd/auth/init/router"
	authDelivery "github.com/yarikTri/archipelago-notes-api/internal/pkg/auth/delivery/http"
	usersRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/auth/repository/postgresql"
	sessionsRepository "github.com/yarikTri/archipelago-notes-api/internal/pkg/auth/repository/redis"
	authUsecase "github.com/yarikTri/archipelago-notes-api/internal/pkg/auth/usecase"
)

func Init(logger logger.Logger) (http.Handler, error) {
	postgresqlDB, err := postgresql.InitPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("error while connecting to postgresql: %v", err)
	}

	redisDB, err := redis.InitRedisDB()
	if err != nil {
		return nil, fmt.Errorf("error while connecting to redis: %v", err)
	}

	usersRepo := usersRepository.NewUsersRepository(postgresqlDB)
	sessionsRepo := sessionsRepository.NewSessionsRepository(redisDB)

	authUsecase := authUsecase.NewUsecase(sessionsRepo, usersRepo)

	authDelivery := authDelivery.NewHandler(authUsecase, logger)

	return router.InitRoutes(authDelivery), nil
}
