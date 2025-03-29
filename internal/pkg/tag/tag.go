package tag

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error)
	UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error
	UpdateTag(ID uuid.UUID, name string) (*models.Tag, error)
	UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error)
	GetNotesByTag(tagID uuid.UUID) ([]models.Note, error)
	GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error)
	LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	GetLinkedTags(tagID uuid.UUID) ([]models.Tag, error)
}

type Repository interface {
	CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error)
	UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error
	UpdateTag(ID uuid.UUID, name string) (*models.Tag, error)
	UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error)
	GetNotesByTag(tagID uuid.UUID) ([]models.Note, error)
	GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error)
	LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	GetLinkedTags(tagID uuid.UUID) ([]models.Tag, error)
}
