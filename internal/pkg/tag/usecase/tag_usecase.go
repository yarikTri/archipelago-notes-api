package usecase

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/errors"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/usecase/dependencies"
)

// Usecase implements tag.Usecase
type Usecase struct {
	tagRepo       tag.TagRepository
	suggesterRepo tag.TagSuggesterRepository
	tagsGraph     dependencies.TagsGraph
}

func NewUsecase(
	tr tag.TagRepository,
	sr tag.TagSuggesterRepository,
	tg dependencies.TagsGraph,
) *Usecase {
	return &Usecase{
		tagRepo:       tr,
		suggesterRepo: sr,
		tagsGraph:     tg,
	}
}

func (u *Usecase) CreateAndLinkTag(name string, noteID, userID uuid.UUID) (*models.Tag, error) {
	tag, err := u.tagRepo.CreateAndLinkTag(name, noteID, userID)
	if err != nil {
		return nil, err
	}

	// Closing eyes on goroutines leakage here.
	go func() {
		if err := u.tagsGraph.UpdateOrCreateTag(tag); err != nil {
			fmt.Printf("Failed to update or create tag: %v\n", err)
		}
	}()

	return tag, nil
}

func (u *Usecase) LinkTagToNote(tagID uuid.UUID, noteID uuid.UUID) error {
	return u.tagRepo.LinkTagToNote(tagID, noteID)
}

func (u *Usecase) UpdateTag(ID uuid.UUID, name string, userID uuid.UUID) (*models.Tag, error) {
	if name == "" {
		return nil, &errors.TagNameEmptyError{}
	}

	return u.tagRepo.UpdateTag(ID, name, userID)
}

func (u *Usecase) UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error {
	return u.tagRepo.UnlinkTagFromNote(tagID, noteID)
}

// func (u *Usecase) UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error) {
// 	if newName == "" {
// 		return nil, &errors.TagNameEmptyError{}
// 	}

// 	return u.tagRepo.UpdateTagForNote(tagID, noteID, newName)
// }

func (u *Usecase) GetNotesByTag(tagID uuid.UUID) ([]models.Note, error) {
	return u.tagRepo.GetNotesByTag(tagID)
}

func (u *Usecase) GetTagsByNoteForUser(noteID, userID uuid.UUID) ([]models.Tag, error) {
	return u.tagRepo.GetTagsByNoteForUser(noteID, userID)
}

func (u *Usecase) LinkTags(tag1ID, tag2ID uuid.UUID) error {
	return u.tagRepo.LinkTags(tag1ID, tag2ID)
}

func (u *Usecase) UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return u.tagRepo.UnlinkTags(tag1ID, tag2ID)
}

func (u *Usecase) GetLinkedTagsForUser(tagID, userID uuid.UUID) ([]models.Tag, error) {
	return u.tagRepo.GetLinkedTagsForUser(tagID, userID)
}

func (u *Usecase) DeleteTag(tagID uuid.UUID) error {
	return u.tagRepo.DeleteTag(tagID)
}

func (u *Usecase) SuggestTags(text string, tagsNum *int) ([]string, error) {
	return u.suggesterRepo.SuggestTags(text, tagsNum)
}

func (u *Usecase) IsTagUsers(userID uuid.UUID, tagID uuid.UUID) (bool, error) {
	tag, err := u.tagRepo.GetTagByID(tagID)
	if err != nil {
		return false, err
	}
	return tag.UserID == userID, nil
}
