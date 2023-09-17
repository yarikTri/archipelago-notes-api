package station

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	GetByID(stationID uint32) (*models.Station, error)
	List() ([]models.Station, error)
	Create(station models.Station) (*models.Station, error)
	DeleteByID(stationID uint32) error
}

type Repository interface {
	GetByID(stationID uint32) (*models.Station, error)
	List() ([]models.Station, error)
	Create(station models.Station) (*models.Station, error)
	DeleteByID(stationID uint32) error
}

type Tables interface {
	Stations() string
	RoutesStations() string
}
