package dirs

import "github.com/yarikTri/archipelago-notes-api/internal/models"

type Usecase interface {
	Get(dirID int) (*models.Dir, error)
	GetTree(dirID int) (*models.DirTree, error)
	Create(name string, parentDirID int) (*models.Dir, error)
	Update(dir *models.Dir) (*models.Dir, error)
	Delete(dirID int) error
}

type Repository interface {
	GetByID(dirID int) (*models.Dir, error)
	GetSubTreeDirsByID(dirID int) ([]*models.Dir, error)
	Create(parentDirID int, name string) (*models.Dir, error)
	Update(dir *models.Dir) (*models.Dir, error)
	DeleteByID(dirID int) error
}
