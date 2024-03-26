package models

import (
	"github.com/gofrs/uuid/v5"
)

type Note struct {
	ID        uuid.UUID
	Title     string
	PlainText string
}

func (n *Note) ToTransfer() NoteTransfer {
	return NoteTransfer{
		ID:        n.ID,
		Title:     n.Title,
		PlainText: n.PlainText,
	}
}

type NoteTransfer struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	PlainText string    `json:"plain_text"`
}
