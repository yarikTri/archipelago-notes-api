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

func (p *PostgreSQL) GetByID(stationID uint32) (*models.Station, error) {
	return nil, nil
}

func (p *PostgreSQL) List() ([]models.Station, error) {
	return nil, nil
}

func (p *PostgreSQL) Create(station models.Station) (*models.Station, error) {
	return nil, nil
}

func (p *PostgreSQL) DeleteByID(stationID uint32) error {
	return nil
}
