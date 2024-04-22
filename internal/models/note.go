package models

import (
	"github.com/gofrs/uuid/v5"
)

type NoteAccess uint8

const (
	UndefinedNoteAccess NoteAccess = iota
	EmptyNoteAccess
	ReadNoteAccess
	WriteNoteAccess
	ModifyNoteAccess
	ManageAccessNoteAccess
)

func NoteAccessFromString(access string) NoteAccess {
	switch access {
	case "e":
		return EmptyNoteAccess
	case "r":
		return ReadNoteAccess
	case "w":
		return WriteNoteAccess
	case "m":
		return ModifyNoteAccess
	case "ma":
		return ManageAccessNoteAccess
	}

	return UndefinedNoteAccess
}

func (na *NoteAccess) String() string {
	switch *na {
	case EmptyNoteAccess:
		return "e"
	case ReadNoteAccess:
		return "r"
	case WriteNoteAccess:
		return "w"
	case ModifyNoteAccess:
		return "m"
	case ManageAccessNoteAccess:
		return "ma"
	}

	return ""
}

type Note struct {
	ID            uuid.UUID `db:"id"`
	DirID         int       `db:"dir_id"`
	Title         string    `db:"title"`
	AutomergeURL  string    `db:"automerge_url"`
	CreatorID     uuid.UUID `db:"creator_id"`
	DefaultAccess string    `db:"default_access"`
	Access        *string   `db:"access"`
}

func (n *Note) ToTransfer(allowedMethods []string) *NoteTransfer {
	return &NoteTransfer{
		ID:            n.ID.String(),
		DirID:         n.DirID,
		AutomergeURL:  n.AutomergeURL,
		Title:         n.Title,
		CreatorID:     n.CreatorID.String(),
		DefaultAccess: n.DefaultAccess,

		AllowedMethods: allowedMethods,
	}
}

type NoteTransfer struct {
	ID            string `json:"id"`
	DirID         int    `json:"dir_id"`
	Title         string `json:"title"`
	AutomergeURL  string `json:"automerge_url"`
	CreatorID     string `json:"creator_id"`
	DefaultAccess string `json:"default_access"`

	AllowedMethods []string `json:"allowed_methods"`
}
