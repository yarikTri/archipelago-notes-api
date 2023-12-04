package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv" // load environment

	flog "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	app "github.com/yarikTri/web-transport-cards/cmd/api/init"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/config"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/postgresql"
	"github.com/yarikTri/web-transport-cards/cmd/api/init/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqIdGetterMock := func(context.Context) (uint32, error) { return 0, nil }
	flogger, err := flog.NewFLogger(reqIdGetterMock)
	if err != nil {
		log.Fatalf("logger can not be defined: %v\n", err)
	}

	db, _, err := postgresql.InitPostgresDB()
	if err != nil {
		flogger.Errorf("error while connecting to database: %v", err)
		return
	}

	router, err := app.Init(db, flogger)
	if err != nil {
		flogger.Errorf("error while launching routes: %v", err)
		return
	}

	var srv server.Server
	if err := srv.Init(os.Getenv(config.ApiListenParamName), router); err != nil {
		flogger.Errorf("error while launching server: %v", err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			flogger.Errorf("server error: %v", err)
			os.Exit(1)
		}
	}()
	flogger.Info("trying to launch server")

	timer := time.AfterFunc(2*time.Second, func() {
		flogger.Infof("server launched at %s", os.Getenv(config.ApiListenParamName))
	})
	defer timer.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	flogger.Info("server gracefully shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		flogger.Errorf("error while shutting down server: %v", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error while loading environment: %v", err)
	}
}
