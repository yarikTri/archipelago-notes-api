package postgresql

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"
	"gorm.io/gorm"
)

// PostgreSQL implements ticket.Repository
type PostgreSQL struct {
	db     *gorm.DB
	tables ticket.Tables
}

func NewPostgreSQL(db *gorm.DB, t ticket.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) GetByID(ticketID uint32) (models.Ticket, error) {
	return models.Ticket{}, nil
}

func (p *PostgreSQL) List() ([]models.Ticket, error) {
	return nil, nil
}

func (p *PostgreSQL) Create(ticket models.Ticket) (models.Ticket, error) {
	return models.Ticket{}, nil
}

func (p *PostgreSQL) DeleteByID(ticketID uint32) error {
	return nil
}
