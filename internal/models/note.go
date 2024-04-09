package models

import (
	"github.com/gofrs/uuid/v5"
)

type Note struct {
	ID        uuid.UUID `db:"id"`
	AutomergeURL string `db:"automerge_url"`
	Title     string    `db:"title"`
}

func (n *Note) ToTransfer() NoteTransfer {
	return NoteTransfer{
		ID:        n.ID.String(),
		AutomergeURL: n.AutomergeURL,
		Title:     n.Title,
	}
}

type NoteTransfer struct {
	ID        string `json:"id"`
	AutomergeURL string `json:"automerge_url"`
	Title     string `json:"title"`
}
