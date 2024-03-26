package notes

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	GetByID(noteID uuid.UUID) (*models.Note, error)
	List() ([]models.Note, error)
	Create(title string) (*models.Note, error)
	Update(route models.Note) (*models.Note, error)
	DeleteByID(noteID uuid.UUID) error
}

type Repository interface {
	GetByID(nodeID uuid.UUID) (*models.Note, error)
	List() ([]models.Note, error)
	Create(title string) (*models.Note, error)
	Update(route models.Note) (*models.Note, error)
	DeleteByID(nodeID uuid.UUID) error
}
