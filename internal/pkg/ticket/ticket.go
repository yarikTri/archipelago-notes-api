package ticket

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	GetByID(ticketID int) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	FormByID(ticketID int) (models.Ticket, error)
	ApproveByID(ticketID int) (models.Ticket, error)
	RejectByID(ticketID int) (models.Ticket, error)
	DeleteByID(ticketID int) error
	AddRoute(ticketID, routeID int) (models.Ticket, error)
	DeleteRoute(ticketID, routeID int) (models.Ticket, error)
	GetTicketDraftByCreatorID(creatorID int) *models.Ticket
}

type Repository interface {
	GetByID(ticketID int) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	FormByID(ticketID int) (models.Ticket, error)
	ApproveByID(ticketID int) (models.Ticket, error)
	RejectByID(ticketID int) (models.Ticket, error)
	DeleteByID(ticketID int) error
	AddRoute(ticketID, routeID int) (models.Ticket, error)
	DeleteRoute(ticketID, routeID int) (models.Ticket, error)
	GetTicketDraftByCreatorID(creatorID int) *models.Ticket
}

type Tables interface {
	Tickets() string
	Routes() string
	RoutesStations() string
	Stations() string
}
