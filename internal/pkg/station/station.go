package station

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	ListByRoute(routeID uint32) ([]models.Station, error)
}

type Repository interface {
	ListByRoute(routeID uint32) ([]models.Station, error)
}

type Tables interface {
	Stations() string
	RoutesStations() string
}
