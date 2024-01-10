package minio

import (
	"errors"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	Secret    string
}

func initMinioConfig() (MinioConfig, error) {
	cfg := MinioConfig{
		Endpoint:  os.Getenv("MINIO_LISTEN_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		Secret:    os.Getenv("MINIO_SECRET_KEY"),
	}

	if strings.TrimSpace(cfg.Endpoint) == "" ||
		strings.TrimSpace(cfg.AccessKey) == "" ||
		strings.TrimSpace(cfg.Secret) == "" {

		return MinioConfig{}, errors.New("invalid minio config")
	}

	return cfg, nil
}

func MakeS3MinioClient() (*minio.Client, error) {
	cfg, err := initMinioConfig()
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(
		cfg.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.Secret, ""),
		},
	)
	if err != nil {
		return nil, err
	}

	if minioClient.IsOffline() {
		return nil, errors.New("Minio client is offline")
	}

	return minioClient, nil
}
