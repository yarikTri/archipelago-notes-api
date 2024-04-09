package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type CreateNoteRequest struct {
	Title     string `json:"title" valid:"required"`
	AutomergeURL string `json:"automerge_url" valid:"required"`
}

func (cnr *CreateNoteRequest) validate() error {
	_, err := valid.ValidateStruct(cnr)
	return err
}

type UpdateNoteRequest struct {
	ID        string `json:"id"`
	AutomergeURL string `json:"automerge_url" valid:"required"`
	Title     string `json:"title" valid:"required"`
}

func (unr *UpdateNoteRequest) validate() error {
	_, err := valid.ValidateStruct(unr)
	return err
}

func (unr *UpdateNoteRequest) ToNote() models.Note {
	id, _ := uuid.FromString(unr.ID)
	return models.Note{
		ID:        id,
		AutomergeURL: unr.AutomergeURL,
		Title:     unr.Title,
	}
}

type ListNotesResponse struct {
	Notes []models.NoteTransfer `json:"notes"`
}
