package postgresql

import (
	"errors"
	"fmt"
	"time"

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
	var ticket models.Ticket
	if err := p.db.Preload("Routes").First(&ticket, ticketID).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) List() ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := p.db.Preload("Routes").Where("state IN ('?', '?', '?', '?')",
		models.REJECTED_STATE, models.FORMED_STATE, models.APPROVED_STATE, models.ENDED_STATE).Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (p *PostgreSQL) ListAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := p.db.Preload("Routes").Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (p *PostgreSQL) Create(ticket models.Ticket) (models.Ticket, error) {
	if err := p.db.Create(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}
	return ticket, nil
}

func (p *PostgreSQL) FormByID(ticketID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.DRAFT_STATE {
		return models.Ticket{}, errors.New("Invalid ticket's state to form, has to be 'draft'")
	}

	ticket.State = models.FORMED_STATE
	ticket.FormTime = time.Now()
	if err = p.db.Preload("Routes").Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) ApproveByID(ticketID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.FORMED_STATE {
		return models.Ticket{}, errors.New("Invalid ticket's state to approve, has to be 'formed'")
	}

	ticket.State = models.APPROVED_STATE
	ticket.ApproveTime = time.Now()
	ticket.EndTime = time.Now().Add(90 * time.Minute)
	if err = p.db.Preload("Routes").Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) RejectByID(ticketID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.FORMED_STATE {
		return models.Ticket{}, errors.New("Invalid ticket's state to reject, has to be 'formed'")
	}

	ticket.State = models.REJECTED_STATE
	if err = p.db.Preload("Routes").Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) DeleteByID(ticketID int) error {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return err
	}

	if ticket.State != models.DRAFT_STATE {
		return errors.New(fmt.Sprintf("Invalid ticket's state to delete, has to be %s", models.DRAFT_STATE))
	}

	ticket.State = models.DELETED_STATE
	return p.db.Save(&ticket).Error
}

func (p *PostgreSQL) AddRoute(ticketID, routeID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	var route models.Route
	if err := p.db.First(&route, routeID).Error; err != nil {
		return models.Ticket{}, err
	}

	ticket.Routes = append(ticket.Routes, route)

	if err = p.db.Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) DeleteRoute(ticketID, routeID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	var route models.Route
	if err := p.db.First(&route, routeID).Error; err != nil {
		return models.Ticket{}, err
	}

	currRoutes := ticket.Routes
	for ind, _route := range ticket.Routes {
		if _route.ID == route.ID {
			ticket.Routes = append(currRoutes[:ind], currRoutes[ind+1:]...)
			if err = p.db.Save(&ticket).Error; err != nil {
				return models.Ticket{}, err
			}
			return ticket, nil
		}
	}
	return ticket, nil
}

func (p *PostgreSQL) GetTicketDraftByCreatorID(creatorID int) *models.Ticket {
	var ticket models.Ticket
	p.db.Where("state = ? AND creator_id = ?", models.DRAFT_STATE, creatorID).First(&ticket)

	return &ticket
}
