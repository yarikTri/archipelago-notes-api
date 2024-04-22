package usecase

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/users"
)

// Usecase implements users.Usecase
type Usecase struct {
	repo users.Repository
}

func NewUsecase(ur users.Repository) *Usecase {
	return &Usecase{
		repo: ur,
	}
}

func (u *Usecase) GetByID(userID uuid.UUID) (*models.User, error) {
	return u.repo.GetByID(userID)
}

func (u *Usecase) Search(query string) ([]*models.User, error) {
	return u.repo.Search(query)
}

func (u *Usecase) SetRootDirByID(userID uuid.UUID, dirID int) error {
	return u.repo.SetRootDirByID(userID, dirID)
}
