package usecase

import (
	"github.com/yarikTri/archipelago-nodes-api/internal/models"
)

// Usecase implements notes.Usecase
type Usecase struct {
	repo notes.Repository
}

func NewUsecase(rr notes.Repository) *Usecase {
	return &Usecase{
		repo: rr,
	}
}

func (u *Usecase) GetByID(routeID int) (models.Route, error) {
	return u.repo.GetByID(routeID)
}

func (u *Usecase) List() ([]models.Route, error) {
	return u.repo.List()
}

func (u *Usecase) Create(route models.Route) (models.Route, error) {
	return u.repo.Create(route)
}

func (u *Usecase) Update(route models.Route) (models.Route, error) {
	return u.repo.Update(route)
}

func (u *Usecase) DeleteByID(routeID int) error {
	return u.repo.DeleteByID(routeID)
}
