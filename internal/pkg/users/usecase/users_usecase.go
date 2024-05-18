package usecase

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/clients/invitations/email"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/users"
)

// Usecase implements users.Usecase
type Usecase struct {
	repo                    users.Repository
	emailConfirmationClient email.IEmailConfirmationClient
}

func NewUsecase(ur users.Repository, ecc email.IEmailConfirmationClient) *Usecase {
	return &Usecase{
		repo:                    ur,
		emailConfirmationClient: ecc,
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

func (u *Usecase) SendEmailConfirmation(userID uuid.UUID) error {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return err
	}

	return u.emailConfirmationClient.SendConfirmation(user.Email, user.ID.String())
}

func (u *Usecase) ConfirmEmail(userID uuid.UUID) error {
	return u.repo.ConfirmEmail(userID)
}
