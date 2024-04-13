package notes

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	GetByID(noteID uuid.UUID) (*models.Note, error)
	List(userID uuid.UUID) ([]*models.Note, error)
	Create(dirID int, automergeURL, title string, creatorID uuid.UUID) (*models.Note, error)
	Update(note models.Note) (*models.Note, error)
	DeleteByID(noteID uuid.UUID) error

	GetUserAccess(noteID uuid.UUID, userID uuid.UUID) (models.NoteAccess, error)
	SetUserAccess(noteID uuid.UUID, userID uuid.UUID, access models.NoteAccess, sendInvitation bool) error
	CheckOwner(noteID uuid.UUID, userID uuid.UUID) (bool, error)

	AttachNoteToSummary(summID, noteID uuid.UUID) error
	DettachNoteFromSummary(summID, noteID uuid.UUID) error
	GetSummaryListByNote(noteID uuid.UUID) ([]uuid.UUID, []uuid.UUID, error)
}

type Repository interface {
	GetByID(noteID uuid.UUID) (*models.Note, error)
	List(userID uuid.UUID) ([]*models.Note, error)
	Create(dirID int, automergeURL, title string, creatorID uuid.UUID) (*models.Note, error)
	Update(note models.Note) (*models.Note, error)
	DeleteByID(noteID uuid.UUID) error

	GetUserAccess(noteID uuid.UUID, userID uuid.UUID) (models.NoteAccess, error)
	SetUserAccess(noteID uuid.UUID, userID uuid.UUID, access models.NoteAccess) error

	AttachNoteToSummary(summID, noteID uuid.UUID) error
	DettachNoteFromSummary(summID, noteID uuid.UUID) error
	GetSummaryListByNote(noteID uuid.UUID) ([]models.SummaryIDStatus, error)
}
