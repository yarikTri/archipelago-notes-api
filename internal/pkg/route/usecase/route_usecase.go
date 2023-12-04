package usecase

import (
	"github.com/google/uuid"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
)

// Usecase implements route.Usecase
type Usecase struct {
	repo route.Repository
}

func NewUsecase(rr route.Repository) *Usecase {
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

func (u *Usecase) Search(subString string) ([]models.Route, error) {
	return u.repo.Search(subString)
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

func (u *Usecase) UpdateImageUUID(routeID int, imageUUID uuid.UUID) error {
	return u.repo.UpdateImageUUID(routeID, imageUUID)
}
