package postgresql

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxIdleConns = 10
	maxOpenConns = 10
)

// Config includes info about postgres DB we want to connect to
type PostgresConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	DBSSLMode  string
}

// InitConfig inits DB configuration from environment variables
func initPostgresConfig() (PostgresConfig, error) { // TODO CHECK FIELDS
	cfg := PostgresConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if strings.TrimSpace(cfg.DBHost) == "" ||
		strings.TrimSpace(cfg.DBPort) == "" ||
		strings.TrimSpace(cfg.DBUser) == "" ||
		strings.TrimSpace(cfg.DBName) == "" ||
		strings.TrimSpace(cfg.DBPassword) == "" ||
		strings.TrimSpace(cfg.DBSSLMode) == "" {

		fmt.Print(cfg)

		return PostgresConfig{}, errors.New("invalid db config")
	}

	return cfg, nil
}

// NewPostgresDB connects to chosen postgreSQL database
// and returns interaction interface of the database
func InitPostgresDB() (*gorm.DB, PostgreSQLTables, error) {
	cfg, err := initPostgresConfig()
	if err != nil {
		return nil, PostgreSQLTables{}, fmt.Errorf("can't init postgresql: %w", err)
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})
	if err != nil {
		return nil, PostgreSQLTables{}, err
	}

	// db.SetMaxIdleConns(maxIdleConns)
	// db.SetMaxOpenConns(maxOpenConns)

	// err = db.Ping()
	// if err != nil {
	// 	errClose := db.Close()
	// 	if errClose != nil {
	// 		return nil, PostgreSQLTables{},
	// 			fmt.Errorf("can't close postgresql (%w) after failed ping: %w", errClose, err)
	// 	}
	// 	return nil, PostgreSQLTables{}, err
	// }

	return db, PostgreSQLTables{}, nil
}

type PostgreSQLTables struct{}

func (pt PostgreSQLTables) Stations() string {
	return "Stations"
}

func (pt PostgreSQLTables) Routes() string {
	return "Routes"
}

func (pt PostgreSQLTables) RoutesStations() string {
	return "Routes_Stations"
}

func (pt PostgreSQLTables) RoutesTickets() string {
	return "Routes_Tickets"
}

func (pt PostgreSQLTables) Tickets() string {
	return "Tickets"
}
