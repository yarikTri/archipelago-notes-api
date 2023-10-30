package ticket

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	GetByID(ticketID uint32) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	DeleteByID(ticketID uint32) error
}

type Repository interface {
	GetByID(ticketID uint32) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	DeleteByID(ticketID uint32) error
}

type Tables interface {
	Tickets() string
	Routes() string
	RoutesStations() string
	Stations() string
}
