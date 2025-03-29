package http

import (
	valid "github.com/asaskevich/govalidator"
)

// CreateAndLinkTagRequest represents the request body for creating and linking a tag
type CreateAndLinkTagRequest struct {
	Name   string `json:"name" valid:"required"`
	NoteID string `json:"note_id" valid:"required"`
}

// UnlinkTagRequest represents the request body for unlinking a tag from a note
type UnlinkTagRequest struct {
	TagID  string `json:"tag_id" valid:"required"`
	NoteID string `json:"note_id" valid:"required"`
}

// LinkTagsRequest represents the request body for linking two tags
type LinkTagsRequest struct {
	Tag1ID string `json:"tag1_id" valid:"required"`
	Tag2ID string `json:"tag2_id" valid:"required"`
}

// UnlinkTagsRequest represents the request body for unlinking two tags
type UnlinkTagsRequest struct {
	Tag1ID string `json:"tag1_id" valid:"required"`
	Tag2ID string `json:"tag2_id" valid:"required"`
}

type UpdateTagRequest struct {
	TagID string `json:"tag_id" valid:"required"`
	Name  string `json:"name" valid:"required"`
}

func (utr *UpdateTagRequest) validate() error {
	_, err := valid.ValidateStruct(utr)
	return err
}

type UpdateTagForNoteRequest struct {
	TagID  string `json:"tag_id" valid:"required"`
	NoteID string `json:"note_id" valid:"required"`
	Name   string `json:"name" valid:"required"`
}

func (utr *UpdateTagForNoteRequest) validate() error {
	_, err := valid.ValidateStruct(utr)
	return err
}
