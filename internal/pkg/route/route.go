package route

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	GetByID(routeID int) (models.Route, error)
	List() ([]models.Route, error)
	Search(subString string) ([]models.Route, error)
	DeleteByID(routeID int) error
}

type Repository interface {
	GetByID(routeID int) (models.Route, error)
	List() ([]models.Route, error)
	Search(subString string) ([]models.Route, error)
	DeleteByID(routeID int) error
}

type Tables interface {
	Routes() string
	RoutesStations() string
	Stations() string
}
