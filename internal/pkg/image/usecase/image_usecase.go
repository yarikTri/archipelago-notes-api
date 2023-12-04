package usecase

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/yarikTri/web-transport-cards/internal/pkg/image"
)

// Usecase implements route.Usecase
type Usecase struct {
	repo image.Repository
}

func NewUsecase(ir image.Repository) *Usecase {
	return &Usecase{
		repo: ir,
	}
}

func (u *Usecase) Put(ctx context.Context, image io.Reader, size int64) (uuid.UUID, error) {
	return u.repo.Put(ctx, image, size)
}

func (u *Usecase) Delete(ctx context.Context, uuid uuid.UUID) error {
	return u.repo.Delete(ctx, uuid)
}
