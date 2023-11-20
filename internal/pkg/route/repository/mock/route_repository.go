package mock

import (
	"fmt"
	"strings"

	"github.com/yarikTri/web-transport-cards/cmd/api/init/db/mock"
	"github.com/yarikTri/web-transport-cards/internal/models"
)

// Mock implements route.Repository
type Mock struct {
	db mock.MockDB
}

func NewMock(db mock.MockDB) *Mock {
	return &Mock{
		db: db,
	}
}

func (m *Mock) GetByID(routeID int) (models.Route, error) {
	return m.db.Routes[fmt.Sprint(routeID)], nil
}

func (m *Mock) List() ([]models.Route, error) {
	routes := make([]models.Route, 0)

	for _, route := range m.db.Routes {
		routes = append(routes, route)
	}

	return routes, nil
}

func (m *Mock) Search(subRoute string) ([]models.Route, error) {
	routes := make([]models.Route, 0)

	for _, route := range m.db.Routes {
		if strings.Contains(strings.ToLower(route.Name), strings.ToLower(subRoute)) {
			routes = append(routes, route)
		}
	}

	return routes, nil
}

func (m *Mock) DeleteByID(routeID int) error {
	delete(m.db.Routes, fmt.Sprint(routeID))

	return nil
}
