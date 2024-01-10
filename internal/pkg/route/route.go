package route

import (
	"github.com/google/uuid"
	"github.com/yarikTri/web-transport-cards/internal/models"
)

type Usecase interface {
	GetByID(routeID int) (models.Route, error)
	List() ([]models.Route, error)
	Search(subString string) ([]models.Route, error)
	Create(route models.Route) (models.Route, error)
	Update(route models.Route) (models.Route, error)
	DeleteByID(routeID int) error
	UpdateImageUUID(routeID int, imageUUID uuid.UUID) error
}

type Repository interface {
	GetByID(routeID int) (models.Route, error)
	List() ([]models.Route, error)
	Search(subString string) ([]models.Route, error)
	Create(route models.Route) (models.Route, error)
	Update(route models.Route) (models.Route, error)
	DeleteByID(routeID int) error
	UpdateImageUUID(routeID int, imageUUID uuid.UUID) error
}
