package ticket

import "github.com/yarikTri/web-transport-cards/internal/models"

type Usecase interface {
	GetByID(ticketID int) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	FormDraft(ticketID int) (models.Ticket, error)
	ApproveByID(ticketID, moderatorID int) (models.Ticket, error)
	RejectByID(ticketID, moderatorID int) (models.Ticket, error)
	DeleteByID(ticketID int) error

	GetDraft(creatorID int) (models.Ticket, error)
	DeleteDraft(creatorID int) error
	AddRoute(creatorID, routeID int) (models.Ticket, error)
	DeleteRoute(creatorID, routeID int) (models.Ticket, error)

	FinalizeWriting(ticketID int) (models.Ticket, error)
}

type Repository interface {
	GetByID(ticketID int) (models.Ticket, error)
	List() ([]models.Ticket, error)
	Create(ticket models.Ticket) (models.Ticket, error)
	DeleteByID(ticketID int) error
	ApproveByID(ticketID, moderatorID int) (models.Ticket, error)
	RejectByID(ticketID, moderatorID int) (models.Ticket, error)

	FinalizeWriting(ticketID int) (models.Ticket, error)
}

type DraftRepository interface {
	GetTicketDraft(creatorID int) (models.Ticket, error)
	SetTicketDraft(ticket models.Ticket) (models.Ticket, error)
	DelTicketDraft(creatorID int) error

	AddRoute(creatorID int, route models.Route) (models.Ticket, error)
	DeleteRoute(creatorID int, routeID int) (models.Ticket, error)
}
