package models

import (
	"github.com/gofrs/uuid/v5"
)

type Note struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	PlainText string    `db:"plain_text"`
}

func (n *Note) ToTransfer() NoteTransfer {
	return NoteTransfer{
		ID:        n.ID.String(),
		Title:     n.Title,
		PlainText: n.PlainText,
	}
}

type NoteTransfer struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	PlainText string `json:"plain_text"`
}
