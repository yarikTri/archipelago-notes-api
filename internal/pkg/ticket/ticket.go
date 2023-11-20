package ticket

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	GetByID(ticketID int) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	DeleteByID(ticketID int) error
}

type Repository interface {
	GetByID(ticketID int) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	DeleteByID(ticketID int) error
}

type Tables interface {
	Tickets() string
	Routes() string
	RoutesStations() string
	Stations() string
}
