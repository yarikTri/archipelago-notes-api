package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv" // load environment

	flog "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	app "github.com/yarikTri/archipelago-notes-api/cmd/api/init"
	"github.com/yarikTri/archipelago-notes-api/cmd/api/init/config"
	"github.com/yarikTri/archipelago-notes-api/cmd/api/init/server"
	"github.com/yarikTri/archipelago-notes-api/cmd/common/init/db/postgresql"
)

// @title		Archipelago Notes API
// @version		1.0.1
// @description	Notes API

// @contact.name   Yaroslav Kuzmin
// @contact.email  yarik1448kuzmin@gmail.com

// @host localhost:8080
// @schemes https http
// @BasePath /

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqIdGetterMock := func(context.Context) (uint32, error) { return 0, nil }
	flogger, err := flog.NewFLogger(reqIdGetterMock)
	if err != nil {
		log.Fatalf("logger can not be defined: %v\n", err)
	}

	db, err := postgresql.InitPostgresDB()
	if err != nil {
		flogger.Errorf("error while connecting to database: %v", err)
		return
	}

	openAiUrl := os.Getenv(config.OpenAIUrlParamName)
	if openAiUrl == "" {
		flogger.Errorf("OPENAI_URL is not set")
		return
	}

	tagSuggesterModel := os.Getenv(config.TagSuggesterModelParamName)
	if tagSuggesterModel == "" {
		flogger.Errorf("TAG_SUGGESTER_MODEL is not set")
		return
	}

	defaultGenerateTagNum, err := strconv.Atoi(os.Getenv(config.DefaultGenerateTagNumParamName))
	if err != nil {
		flogger.Errorf("error while converting %s to int: %v", config.DefaultGenerateTagNumParamName, err)
		return
	}

	qdrantHost := os.Getenv(config.QdrantHostParamName)
	if qdrantHost == "" {
		flogger.Errorf("%s is not set", config.QdrantHostParamName)
		return
	}

	qdrantPort := os.Getenv(config.QdrantPortParamName)
	if qdrantPort == "" {
		flogger.Errorf("%s is not set", config.QdrantPortParamName)
		return
	}

	router, err := app.Init(db, flogger, openAiUrl, tagSuggesterModel, defaultGenerateTagNum, qdrantHost, qdrantPort)
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

	timer := time.AfterFunc(1*time.Second, func() {
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
