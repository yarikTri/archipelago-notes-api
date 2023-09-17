package usecase

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/station"
)

// Usecase implements station.Usecase
type Usecase struct {
	repo station.Repository
}

func NewUsecase(sr station.Repository) *Usecase {
	return &Usecase{
		repo: sr,
	}
}

func (u *Usecase) GetByID(stationID uint32) (*models.Station, error) {
	return u.repo.GetByID(stationID)
}

func (u *Usecase) List() ([]models.Station, error) {
	return u.repo.List()
}

func (u *Usecase) Create(station models.Station) (*models.Station, error) {
	return u.repo.Create(station)
}

func (u *Usecase) DeleteByID(stationID uint32) error {
	return u.DeleteByID(stationID)
}
