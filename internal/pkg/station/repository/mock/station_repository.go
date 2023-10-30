package mock

import (
	"fmt"

	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/mock"
	"github.com/yarikTri/web-transport-cards/internal/models"
)

// Mock implements station.Repository
type Mock struct {
	db mock.MockDB
}

func NewMock(db mock.MockDB) *Mock {
	return &Mock{
		db: db,
	}
}

func (m *Mock) ListByRoute(routeID uint32) ([]models.Station, error) {
	stations := make([]models.Station, 0)

	stationsIDs := m.db.RoutesStation[fmt.Sprint(routeID)]

	for _, stationID := range stationsIDs {
		stations = append(stations, m.db.Stations[stationID])
	}

	return stations, nil
}
