package postgresql

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
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
		DBHost:     os.Getenv("POSTGRESQL_HOST"),
		DBPort:     os.Getenv("POSTGRESQL_PORT"),
		DBUser:     os.Getenv("POSTGRESQL_USER"),
		DBName:     os.Getenv("POSTGRESQL_NAME"),
		DBPassword: os.Getenv("POSTGRESQL_PASSWORD"),
		DBSSLMode:  os.Getenv("POSTGRESQL_SSLMODE"),
	}

	if strings.TrimSpace(cfg.DBHost) == "" ||
		strings.TrimSpace(cfg.DBPort) == "" ||
		strings.TrimSpace(cfg.DBUser) == "" ||
		strings.TrimSpace(cfg.DBName) == "" ||
		strings.TrimSpace(cfg.DBPassword) == "" ||
		strings.TrimSpace(cfg.DBSSLMode) == "" {

		fmt.Print(cfg)

		return PostgresConfig{}, errors.New("invalid postgresql config")
	}

	return cfg, nil
}

// NewPostgresDB connects to chosen postgreSQL database
// and returns interaction interface of the database
func InitPostgresDB() (*sqlx.DB, error) {
	cfg, err := initPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("can't init postgresql: %w", err)
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	err = db.Ping()
	if err != nil {
		errClose := db.Close()
		if errClose != nil {
			return nil, fmt.Errorf("can't close postgresql (%w) after failed ping: %w", errClose, err)
		}
		return nil, err
	}

	return db, nil
}
