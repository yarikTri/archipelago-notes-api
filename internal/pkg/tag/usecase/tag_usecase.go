package usecase

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/errors"
)

// Usecase implements tag.Usecase
type Usecase struct {
	tagRepo       tag.TagRepository
	suggesterRepo tag.TagSuggesterRepository
}

func NewUsecase(tr tag.TagRepository, sr tag.TagSuggesterRepository) *Usecase {
	return &Usecase{
		tagRepo:       tr,
		suggesterRepo: sr,
	}
}

func (u *Usecase) CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error) {
	return u.tagRepo.CreateAndLinkTag(name, noteID)
}

func (u *Usecase) LinkExistingTag(tagID uuid.UUID, noteID uuid.UUID) error {
	return u.tagRepo.LinkExistingTag(tagID, noteID)
}

func (u *Usecase) UpdateTag(ID uuid.UUID, name string) (*models.Tag, error) {
	if name == "" {
		return nil, &errors.TagNameEmptyError{}
	}

	return u.tagRepo.UpdateTag(ID, name)
}

func (u *Usecase) UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error {
	return u.tagRepo.UnlinkTagFromNote(tagID, noteID)
}

func (u *Usecase) UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error) {
	if newName == "" {
		return nil, &errors.TagNameEmptyError{}
	}

	return u.tagRepo.UpdateTagForNote(tagID, noteID, newName)
}

func (u *Usecase) GetNotesByTag(tagID uuid.UUID) ([]models.Note, error) {
	return u.tagRepo.GetNotesByTag(tagID)
}

func (u *Usecase) GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error) {
	return u.tagRepo.GetTagsByNote(noteID)
}

func (u *Usecase) LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return u.tagRepo.LinkTags(tag1ID, tag2ID)
}

func (u *Usecase) UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return u.tagRepo.UnlinkTags(tag1ID, tag2ID)
}

func (u *Usecase) GetLinkedTags(tagID uuid.UUID) ([]models.Tag, error) {
	return u.tagRepo.GetLinkedTags(tagID)
}

func (u *Usecase) DeleteTag(tagID uuid.UUID) error {
	return u.tagRepo.DeleteTag(tagID)
}

func (u *Usecase) SuggestTags(text string) ([]string, error) {
	return u.suggesterRepo.SuggestTags(text)
}
