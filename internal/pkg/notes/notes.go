package notes

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	GetByID(noteID uuid.UUID) (*models.Note, error)
	List() ([]*models.Note, error)
	Create(dirID int, automergeURL, title string) (*models.Note, error)
	Update(note models.Note) (*models.Note, error)
	DeleteByID(noteID uuid.UUID) error
}

type Repository interface {
	GetByID(noteID uuid.UUID) (*models.Note, error)
	List() ([]*models.Note, error)
	ListByDirIds(dirIDs []int) ([]*models.Note, error)
	Create(dirID int, automergeURL, title string) (*models.Note, error)
	Update(note models.Note) (*models.Note, error)
	DeleteByID(noteID uuid.UUID) error
}
