package dependencies

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

// TODO: top of mind name variant, rename
type TagsGraph interface {
	UpdateOrCreateTag(tag *models.Tag) error
	ListClosestTagsIds(tag *models.Tag, limit uint32) ([]uuid.UUID, error)
	DeleteByID(tagID uuid.UUID) error
}
