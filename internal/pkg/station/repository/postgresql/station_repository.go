package postgresql

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/station"
	"gorm.io/gorm"
)

// PostgreSQL implements route.Repository
type PostgreSQL struct {
	db     *gorm.DB
	tables station.Tables
}

func NewPostgreSQL(db *gorm.DB, t station.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) ListByRoute(routeID uint32) ([]models.Station, error) {
	return nil, nil
}
