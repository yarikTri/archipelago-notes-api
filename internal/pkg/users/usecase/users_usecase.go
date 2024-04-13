package usecase

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/notes"
)

// Usecase implements notes.Usecase
type Usecase struct {
	repo notes.Repository
}

func NewUsecase(rr notes.Repository) *Usecase {
	return &Usecase{
		repo: rr,
	}
}

func (u *Usecase) GetByID(noteID uuid.UUID) (*models.Note, error) {
	return u.repo.GetByID(noteID)
}

func (u *Usecase) List() ([]*models.Note, error) {
	return u.repo.List()
}

func (u *Usecase) Create(dirID int, automergeURL, title string) (*models.Note, error) {
	return u.repo.Create(dirID, automergeURL, title)
}

func (u *Usecase) Update(note models.Note) (*models.Note, error) {
	return u.repo.Update(note)
}

func (u *Usecase) DeleteByID(noteID uuid.UUID) error {
	return u.repo.DeleteByID(noteID)
}
