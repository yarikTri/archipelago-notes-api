package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-nodes-api/internal/models"
)

type CreateNoteRequest struct {
	Title     string `json:"title" valid:"required"`
	PlainText string `json:"plain_text"`
}

func (cnr *CreateNoteRequest) validate() error {
	_, err := valid.ValidateStruct(cnr)
	return err
}

func (cnr *CreateNoteRequest) ToNode() models.Note {
	return models.Note{
		Title:     cnr.Title,
		PlainText: cnr.PlainText,
	}
}

type UpdateNoteRequest struct {
	ID        string `json:"id"`
	Title     string `json:"title" valid:"required"`
	PlainText string `json:"plain_text,omitempty"`
}

func (unr *UpdateNoteRequest) validate() error {
	_, err := valid.ValidateStruct(unr)
	return err
}

func (unr *UpdateNoteRequest) ToNode() models.Note {
	id, _ := uuid.FromString(unr.ID)
	return models.Note{
		ID:        id,
		Title:     unr.Title,
		PlainText: unr.PlainText,
	}
}

type ListNodesResponse struct {
	Nodes []models.NoteTransfer `json:"notes"`
}
