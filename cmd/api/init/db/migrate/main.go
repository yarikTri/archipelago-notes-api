package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/postgresql"
	"github.com/yarikTri/web-transport-cards/internal/models"
)

func main() {
	db, _, err := postgresql.InitPostgresDB()
	if err != nil {
		log.Fatalf("error while connecting to database: %v", err)
		return
	}

	err = db.AutoMigrate(&models.User{}, &models.Route{}, models.Ticket{})
	if err != nil {
		panic("cant migrate db: " + err.Error())
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error while loading environment: %v", err)
	}
}
