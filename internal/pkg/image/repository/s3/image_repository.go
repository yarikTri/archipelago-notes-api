package s3

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type S3MinioImageStorage struct {
	bucket string
	cl     *minio.Client
}

func NewS3MinioImageStorage(bucket string, client *minio.Client) *S3MinioImageStorage {
	return &S3MinioImageStorage{
		bucket: bucket,
		cl:     client,
	}
}

func (s *S3MinioImageStorage) Put(ctx context.Context, image io.Reader, size int64) (uuid.UUID, error) {
	objUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = s.cl.PutObject(ctx, s.bucket, objUUID.String(),
		image, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return uuid.Nil, err
	}

	return objUUID, nil
}

func (s *S3MinioImageStorage) Delete(ctx context.Context, uuid uuid.UUID) error {
	return s.cl.RemoveObject(ctx, s.bucket, uuid.String(), minio.RemoveObjectOptions{})
}
