package postgresql

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
	"gorm.io/gorm"
)

// PostgreSQL implements route.Repository
type PostgreSQL struct {
	db     *gorm.DB
	tables route.Tables
}

func NewPostgreSQL(db *gorm.DB, t route.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) GetByID(routeID uint32) (models.Route, error) {
	return models.Route{}, nil
}

func (p *PostgreSQL) List() ([]models.Route, error) {
	return nil, nil
}

func (p *PostgreSQL) Create(route models.Route) (models.Route, error) {
	return models.Route{}, nil
}

func (p *PostgreSQL) DeleteByID(routeID uint32) error {
	return nil
}
