package tag

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error)
	UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error
	GetTag(ID uuid.UUID) (*models.Tag, error)
	GetAllTags() ([]models.Tag, error)
	UpdateTag(ID uuid.UUID, name string) error
	DeleteTag(ID uuid.UUID) error
	UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) error
	GetNotesByTag(tagID uuid.UUID) ([]models.Note, error)
	GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error)
}

type Repository interface {
	CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error)
	UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error
	GetTag(ID uuid.UUID) (*models.Tag, error)
	GetAllTags() ([]models.Tag, error)
	UpdateTag(ID uuid.UUID, name string) error
	DeleteTag(ID uuid.UUID) error
	UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) error
	GetNotesByTag(tagID uuid.UUID) ([]models.Note, error)
	GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error)
}
