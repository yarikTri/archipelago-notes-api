package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
)

type redisConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
}

func initRedisConfig() (redisConfig, error) {
	cfg := redisConfig{
		DBHost:     os.Getenv("REDIS_HOST"),
		DBPort:     os.Getenv("REDIS_PORT"),
		DBUser:     os.Getenv("REDIS_USER"),
		DBPassword: os.Getenv("REDIS_PASSWORD"),
	}

	if strings.TrimSpace(cfg.DBHost) == "" ||
		strings.TrimSpace(cfg.DBPort) == "" {

		return redisConfig{}, errors.New("invalid redis config")
	}

	return cfg, nil
}

func InitRedisDB() (*redis.Client, error) {
	cfg, err := initRedisConfig()
	if err != nil {
		return &redis.Client{}, err
	}

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
			Password: cfg.DBPassword,
			Username: cfg.DBUser,
		},
	)

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return &redis.Client{}, err
	}

	return redisClient, nil
}
