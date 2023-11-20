package postgresql

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
	"gorm.io/gorm"
)

// PostgreSQL implements ticket.Repository
type PostgreSQL struct {
	db *gorm.DB
}

func NewPostgreSQL(db *gorm.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) GetByID(ticketID int) (models.Ticket, error) {
	return models.Ticket{}, nil
}

func (p *PostgreSQL) List() ([]models.Ticket, error) {
	return nil, nil
}

func (p *PostgreSQL) Create(ticket models.Ticket) (models.Ticket, error) {
	return models.Ticket{}, nil
}

func (p *PostgreSQL) DeleteByID(ticketID int) error {
	return nil
}
