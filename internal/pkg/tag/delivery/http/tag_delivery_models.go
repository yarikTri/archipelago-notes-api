package http

import (
	valid "github.com/asaskevich/govalidator"
)

type UpdateTagRequest struct {
	ID   string `json:"id" valid:"required"`
	Name string `json:"name" valid:"required"`
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
