package postgresql

import (
	"github.com/google/uuid"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"gorm.io/gorm"
)

// PostgreSQL implements route.Repository
type PostgreSQL struct {
	db *gorm.DB
}

func NewPostgreSQL(db *gorm.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) GetByID(routeID int) (models.Route, error) {
	var route models.Route
	if err := p.db.First(&route, routeID).Error; err != nil {
		return models.Route{}, err
	}

	return route, nil
}

func (p *PostgreSQL) List() ([]models.Route, error) {
	var routes []models.Route
	if err := p.db.Where("active = true").Find(&routes).Error; err != nil {
		return nil, err
	}

	return routes, nil
}

func (p *PostgreSQL) Search(subString string) ([]models.Route, error) {
	var routes []models.Route
	likeStatement := "%" + subString + "%"
	if err := p.db.Where("active = true AND name like ?", likeStatement).Find(&routes).Error; err != nil {
		return nil, err
	}

	return routes, nil
}

func (p *PostgreSQL) Create(route models.Route) (models.Route, error) {
	if err := p.db.Create(&route).Error; err != nil {
		return models.Route{}, err
	}

	return route, nil
}

func (p *PostgreSQL) Update(route models.Route) (models.Route, error) {
	if err := p.db.Save(&route).Error; err != nil {
		return models.Route{}, err
	}

	return route, nil
}

func (p *PostgreSQL) DeleteByID(routeID int) error {
	route, err := p.GetByID(routeID)
	if err != nil {
		return err
	}

	route.Active = false
	return p.db.Save(&route).Error
}

func (p *PostgreSQL) UpdateImageUUID(routeID int, imageUUID uuid.UUID) error {
	route, err := p.GetByID(routeID)
	if err != nil {
		return err
	}

	route.ImageUUID = imageUUID
	return p.db.Save(&route).Error
}
