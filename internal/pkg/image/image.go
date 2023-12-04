package image

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type Usecase interface {
	Put(ctx context.Context, image io.Reader, size int64) (uuid.UUID, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type Repository interface {
	Put(ctx context.Context, image io.Reader, size int64) (uuid.UUID, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}
