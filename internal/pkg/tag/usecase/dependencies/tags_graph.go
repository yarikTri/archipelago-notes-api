package dependencies

import (
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

// TODO: top of mind name variant, rename
type TagsGraph interface {
	UpdateOrCreateTag(tag *models.Tag) error
	ListClosestTags(tag *models.Tag, limit uint32) ([]*models.Tag, error)
}
