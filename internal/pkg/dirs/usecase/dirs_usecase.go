package usecase

import (
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/notes"
)

// Usecase implements notes.Usecase
type Usecase struct {
	dirsRepo  dirs.Repository
	notesRepo notes.Repository
}

func NewUsecase(dr dirs.Repository, nr notes.Repository) *Usecase {
	return &Usecase{
		dirsRepo:  dr,
		notesRepo: nr,
	}
}

func (u *Usecase) Get(dirID int) (*models.Dir, error) {
	return u.dirsRepo.GetByID(dirID)
}

func (u *Usecase) GetTree(rootID int) (*models.DirTree, error) {
	dirs, err := u.dirsRepo.GetSubTreeDirsByID(rootID)
	if err != nil {
		return nil, err
	}

	var dirIDs = make([]int, 0)
	for _, dir := range dirs {
		dirIDs = append(dirIDs, dir.ID)
	}

	notes, err := u.notesRepo.ListByDirIds(dirIDs)
	if err != nil {
		return nil, err
	}

	return models.ToTree(rootID, dirs, notes), nil
}

func (u *Usecase) Create(name string, parentDirID int) (*models.Dir, error) {
	return u.dirsRepo.Create(parentDirID, name)
}

func (u *Usecase) Update(dir *models.Dir) (*models.Dir, error) {
	return u.dirsRepo.Update(dir)
}

func (u *Usecase) Delete(dirID int) error {
	return u.dirsRepo.DeleteByID(dirID)
}
