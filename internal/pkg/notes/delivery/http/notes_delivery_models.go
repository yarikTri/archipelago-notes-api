package http

import (
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type CreateNoteRequest struct {
	DirID        int    `json:"dir_id" valid:"required"`
	Title        string `json:"title" valid:"required"`
	AutomergeURL string `json:"automerge_url" valid:"required"`
}

func (cnr *CreateNoteRequest) validate() error {
	_, err := valid.ValidateStruct(cnr)
	return err
}

type UpdateNoteRequest struct {
	ID            string `json:"id" valid:"required"`
	DirID         int    `json:"dir_id" valid:"required"`
	AutomergeURL  string `json:"automerge_url" valid:"required"`
	Title         string `json:"title" valid:"required"`
	DefaultAccess string `json:"default_access" valid:"required"`
}

func (unr *UpdateNoteRequest) validate() error {
	defaultAccess := models.NoteAccessFromString(unr.DefaultAccess)
	if defaultAccess == models.UndefinedNoteAccess ||
		defaultAccess == models.ModifyNoteAccess ||
		defaultAccess == models.ManageAccessNoteAccess {
		return errors.New(fmt.Sprintf("Invalid default access: %s", unr.DefaultAccess))
	}

	_, err := valid.ValidateStruct(unr)
	return err
}

func (unr *UpdateNoteRequest) ToNote() models.Note {
	id, _ := uuid.FromString(unr.ID)
	return models.Note{
		ID:            id,
		DirID:         unr.DirID,
		AutomergeURL:  unr.AutomergeURL,
		Title:         unr.Title,
		DefaultAccess: unr.DefaultAccess,
	}
}

type ListNotesResponse struct {
	Notes []*models.NoteTransfer `json:"notes"`
}

type SetAccessRequest struct {
	Access         string `json:"access" valid:"required"`
	WithInvitation bool   `json:"with_invitation" valid:"required"`
}

func (sar *SetAccessRequest) validate() error {
	_, err := valid.ValidateStruct(sar)
	return err
}
