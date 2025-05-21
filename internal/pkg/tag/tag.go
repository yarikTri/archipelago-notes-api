package tag

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	CreateAndLinkTag(name string, noteID, userID uuid.UUID) (*models.Tag, error)
	LinkTagToNote(tagID uuid.UUID, noteID uuid.UUID) error
	UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error
	UpdateTag(ID uuid.UUID, name string, userID uuid.UUID) (*models.Tag, error)
	// UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error)
	GetNotesByTag(tagID uuid.UUID) ([]models.Note, error)
	GetTagsByNoteForUser(noteID, userID uuid.UUID) ([]models.Tag, error)
	LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	GetLinkedTagsForUser(tagID, userID uuid.UUID) ([]models.LinkedTag, error)
	UpdateTagsLinkName(tag1ID, tag2ID, userID uuid.UUID, linkName string) error
	DeleteTag(tagID uuid.UUID) error
	SuggestTags(text string, tagsNum *int) ([]string, error)
	IsTagUsers(userID uuid.UUID, tagID uuid.UUID) (bool, error)
	ListClosestTags(tagName string, userID uuid.UUID, limit uint32) ([]models.Tag, error)
}

type TagRepository interface {
	CreateAndLinkTag(name string, noteID, userID uuid.UUID) (*models.Tag, error)
	LinkTagToNote(tagID uuid.UUID, noteID uuid.UUID) error
	UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) (bool, error)
	UpdateTag(ID uuid.UUID, name string, userID uuid.UUID) (*models.Tag, error)
	// UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error)
	GetNotesByTag(tagID uuid.UUID) ([]models.Note, error)
	GetTagsByNoteForUser(noteID, userID uuid.UUID) ([]models.Tag, error)
	LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error
	GetLinkedTagsForUser(tagID, userID uuid.UUID) ([]models.LinkedTag, error)
	UpdateTagsLinkName(tag1ID, tag2ID, userID uuid.UUID, linkName string) error
	DeleteTag(tagID uuid.UUID) error
	GetTagByID(tagID uuid.UUID) (*models.Tag, error)
}

type TagSuggesterRepository interface {
	SuggestTags(text string, tagsNum *int) ([]string, error)
}
