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

func (p *PostgreSQL) ApproveByID(ticketID, moderatorID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.FORMED_STATE {
		return models.Ticket{}, errors.New("Invalid ticket's state to approve, has to be 'formed'")
	}

	ticket.State = models.APPROVED_STATE
	ticket.ModeratorID = &moderatorID
	ticket.ApproveTime = time.Now()
	ticket.EndTime = time.Now().Add(90 * time.Minute)
	if err = p.db.Preload("Routes").Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) RejectByID(ticketID, moderatorID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.FORMED_STATE {
		return models.Ticket{}, errors.New("Invalid ticket's state to reject, has to be 'formed'")
	}

	ticket.State = models.REJECTED_STATE
	ticket.ModeratorID = &moderatorID
	if err = p.db.Preload("Routes").Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (p *PostgreSQL) EndByID(ticketID, moderatorID int) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.APPROVED_STATE {
		return models.Ticket{}, errors.New("Invalid ticket's state to approve, has to be 'formed'")
	}

	ticket.State = models.ENDED_STATE
	ticket.ModeratorID = &moderatorID
	ticket.EndTime = time.Now()
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

func (p *PostgreSQL) GetTicketDraftByCreatorID(creatorID int) (models.Ticket, error) {
	var ticket models.Ticket
	err := p.db.Where("state = ? AND creator_id = ?", models.DRAFT_STATE, creatorID).First(&ticket).Error

	return ticket, err
}

func (p *PostgreSQL) UpdateWriteState(ticketID int, state string) (models.Ticket, error) {
	ticket, err := p.GetByID(ticketID)
	if err != nil {
		return models.Ticket{}, err
	}

	if ticket.State != models.APPROVED_STATE {
		return models.Ticket{}, errors.New(fmt.Sprintf("Invalid ticket's state to write, has to be %s", models.APPROVED_STATE))
	}

	ticket.WriteState = &state
	if err := p.db.Save(&ticket).Error; err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}
