package mock

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
)

// Mock implements ticket.Repository
type Mock struct {
	db map[int]models.Ticket
}

func NewMock(db map[int]models.Ticket) *Mock {
	return &Mock{
		db: db,
	}
}

func (m *Mock) GetByID(ticketID int) (*models.Ticket, error) {
	return nil, nil
}

func (m *Mock) List() ([]models.Ticket, error) {
	return nil, nil
}

func (m *Mock) Create(ticket models.Ticket) (*models.Ticket, error) {
	return nil, nil
}

func (*Mock) DeleteByID(ticketID int) error {
	return nil
}
