package errors

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// NoteNotFoundError represents an error when a note is not found
type NoteNotFoundError struct {
	ID uuid.UUID
}

func (e *NoteNotFoundError) Error() string {
	return fmt.Sprintf("note not found: %v", e.ID)
}

// TagNotFoundError represents an error when a tag is not found
type TagNotFoundError struct {
	ID uuid.UUID
}

func (e *TagNotFoundError) Error() string {
	return fmt.Sprintf("tag not found: %v", e.ID)
}

// TagNameExistsError represents an error when a tag with the same name already exists
type TagNameExistsError struct {
	Name string
	ID   uuid.UUID
}

func (e *TagNameExistsError) Error() string {
	return fmt.Sprintf("tag with name '%s' already exists (ID: %s)", e.Name, e.ID)
}

// TagNameEmptyError represents an error when a tag name is empty
type TagNameEmptyError struct{}

func (e *TagNameEmptyError) Error() string {
	return "tag name cannot be empty"
}

// TagLinkNotFoundError represents an error when a tag-note link is not found
type TagLinkNotFoundError struct {
	TagID  uuid.UUID
	NoteID uuid.UUID
}

func (e *TagLinkNotFoundError) Error() string {
	return fmt.Sprintf("tag-note link not found: tag %v, note %v", e.TagID, e.NoteID)
}

type TagLinkExistsError struct {
	TagID  uuid.UUID
	NoteID uuid.UUID
}

func (e *TagLinkExistsError) Error() string {
	return fmt.Sprintf("tag-note link already exists: tag %v, note %v", e.TagID, e.NoteID)
}

// TagLinkExistsError represents an error when a tag-to-tag link already exists
type TagToTagLinkExistsError struct {
	Tag1ID uuid.UUID
	Tag2ID uuid.UUID
}

func (e *TagToTagLinkExistsError) Error() string {
	return fmt.Sprintf("tags are already linked: %v and %v", e.Tag1ID, e.Tag2ID)
}

// TagToTagLinkNotFoundError represents an error when a tag-to-tag link is not found
type TagToTagLinkNotFoundError struct {
	Tag1ID uuid.UUID
	Tag2ID uuid.UUID
}

func (e *TagToTagLinkNotFoundError) Error() string {
	return fmt.Sprintf("tag-to-tag link not found: %v and %v", e.Tag1ID, e.Tag2ID)
}
