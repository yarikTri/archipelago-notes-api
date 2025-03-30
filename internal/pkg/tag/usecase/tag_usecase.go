package usecase

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/errors"
)

// Usecase implements tag.Usecase
type Usecase struct {
	repo tag.Repository
}

func NewUsecase(rr tag.Repository) *Usecase {
	return &Usecase{
		repo: rr,
	}
}

func (u *Usecase) CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error) {
	return u.repo.CreateAndLinkTag(name, noteID)
}

func (u *Usecase) LinkExistingTag(tagID uuid.UUID, noteID uuid.UUID) error {
	return u.repo.LinkExistingTag(tagID, noteID)
}

func (u *Usecase) UpdateTag(ID uuid.UUID, name string) (*models.Tag, error) {
	if name == "" {
		return nil, &errors.TagNameEmptyError{}
	}

	return u.repo.UpdateTag(ID, name)
}

func (u *Usecase) UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error {
	return u.repo.UnlinkTagFromNote(tagID, noteID)
}

func (u *Usecase) UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error) {
	if newName == "" {
		return nil, &errors.TagNameEmptyError{}
	}

	return u.repo.UpdateTagForNote(tagID, noteID, newName)
}

func (u *Usecase) GetNotesByTag(tagID uuid.UUID) ([]models.Note, error) {
	return u.repo.GetNotesByTag(tagID)
}

func (u *Usecase) GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error) {
	return u.repo.GetTagsByNote(noteID)
}

func (u *Usecase) LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return u.repo.LinkTags(tag1ID, tag2ID)
}

func (u *Usecase) UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return u.repo.UnlinkTags(tag1ID, tag2ID)
}

func (u *Usecase) GetLinkedTags(tagID uuid.UUID) ([]models.Tag, error) {
	return u.repo.GetLinkedTags(tagID)
}

func (u *Usecase) DeleteTag(tagID uuid.UUID) error {
	return u.repo.DeleteTag(tagID)
}
