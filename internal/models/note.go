package models

import (
	"github.com/gofrs/uuid/v5"
)

type Note struct {
	ID           uuid.UUID `db:"id"`
	DirID        int       `db:"dir_id"`
	Title        string    `db:"title"`
	AutomergeURL string    `db:"automerge_url"`
}

func (n *Note) ToTransfer() *NoteTransfer {
	return &NoteTransfer{
		ID:           n.ID.String(),
		DirID:        n.DirID,
		AutomergeURL: n.AutomergeURL,
		Title:        n.Title,
	}
}

type NoteTransfer struct {
	ID           string `json:"id"`
	DirID        int    `json:"dir_id"`
	Title        string `json:"title"`
	AutomergeURL string `json:"automerge_url"`
}
