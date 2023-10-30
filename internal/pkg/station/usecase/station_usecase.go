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

func (u *Usecase) ListByRoute(routeID uint32) ([]models.Station, error) {
	return u.repo.ListByRoute(routeID)
}
