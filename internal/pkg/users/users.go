package users

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	GetByID(userID uuid.UUID) (*models.User, error)
	Search(query string) ([]*models.User, error)
	SetRootDirByID(userID uuid.UUID, dirID int) error
	SendEmailConfirmation(userID uuid.UUID) error
	ConfirmEmail(userID uuid.UUID) error
}

type Repository interface {
	GetByID(userID uuid.UUID) (*models.User, error)
	Search(query string) ([]*models.User, error)
	SetRootDirByID(userID uuid.UUID, dirID int) error
	ConfirmEmail(userID uuid.UUID) error
}
