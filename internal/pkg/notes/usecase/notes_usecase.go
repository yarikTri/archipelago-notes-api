package usecase

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/clients/invitations/email"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/notes"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/users"
)

// Usecase implements notes.Usecase
type Usecase struct {
	noteRepo              notes.Repository
	userRepo              users.Repository
	emailInvitationClient email.IEmailInvitationClient
}

func NewUsecase(nr notes.Repository, ur users.Repository, eic email.IEmailInvitationClient) *Usecase {
	return &Usecase{
		noteRepo:              nr,
		userRepo:              ur,
		emailInvitationClient: eic,
	}
}

func (u *Usecase) GetByID(noteID uuid.UUID) (*models.Note, error) {
	return u.noteRepo.GetByID(noteID)
}

func (u *Usecase) List(userID uuid.UUID) ([]*models.Note, error) {
	return u.noteRepo.List(userID)
}

func (u *Usecase) Create(dirID int, automergeURL, title string, creatorID uuid.UUID) (*models.Note, error) {
	return u.noteRepo.Create(dirID, automergeURL, title, creatorID)
}

func (u *Usecase) Update(note models.Note) (*models.Note, error) {
	return u.noteRepo.Update(note)
}

func (u *Usecase) DeleteByID(noteID uuid.UUID) error {
	return u.noteRepo.DeleteByID(noteID)
}

func (u *Usecase) GetUserAccess(noteID uuid.UUID, userID uuid.UUID) (models.NoteAccess, error) {
	return u.noteRepo.GetUserAccess(noteID, userID)
}

func (u *Usecase) SetUserAccess(noteID uuid.UUID, userID uuid.UUID, access models.NoteAccess, sendInvitation bool) error {
	if access == models.UndefinedNoteAccess {
		return errors.New(fmt.Sprintf("(usecase) Invalid access %s", access.String()))
	}

	if sendInvitation {
		if err := u.sendEmailInvitation(noteID, userID); err != nil {
			return err
		}
	}

	return u.noteRepo.SetUserAccess(noteID, userID, access)
}

func (u *Usecase) CheckOwner(noteID uuid.UUID, userID uuid.UUID) (bool, error) {
	note, err := u.noteRepo.GetByID(noteID)
	if err != nil {
		return false, err
	}

	return note.CreatorID == userID, nil
}

func (u *Usecase) DettachNoteFromSummary(summID, noteID uuid.UUID) error {
	return u.noteRepo.DettachNoteFromSummary(summID, noteID)
}

func (u *Usecase) AttachNoteToSummary(summID, noteID uuid.UUID) error {
	return u.noteRepo.AttachNoteToSummary(summID, noteID)
}

func (u *Usecase) GetSummaryListByNote(noteID uuid.UUID) ([]uuid.UUID, []uuid.UUID, error) {
	notes, err := u.noteRepo.GetSummaryListByNote(noteID)
	if err != nil {
		return nil, nil, err
	}

	nonActiveNotesIDS := make([]uuid.UUID, 0, len(notes)/2)
	activeNotesIDS := make([]uuid.UUID, 0, len(notes)/2)

	for _, note := range notes {
		if note.Active {
			activeNotesIDS = append(activeNotesIDS, note.ID)
		} else {
			nonActiveNotesIDS = append(nonActiveNotesIDS, note.ID)
		}
	}

	return nonActiveNotesIDS, activeNotesIDS, nil
}

func (u *Usecase) sendEmailInvitation(noteID uuid.UUID, userID uuid.UUID) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	return u.emailInvitationClient.SendInvitation(user.Email, email.NoteInvitationType, noteID.String())
}
