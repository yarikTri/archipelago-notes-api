package usecase

import (
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"
)

// Usecase implements ticket.Usecase
type Usecase struct {
	repo ticket.Repository
}

func NewUsecase(tr ticket.Repository) *Usecase {
	return &Usecase{
		repo: tr,
	}
}

func (u *Usecase) GetByID(ticketID int) (models.Ticket, error) {
	return u.repo.GetByID(ticketID)
}

func (u *Usecase) List() ([]models.Ticket, error) {
	return u.repo.List()
}

func (u *Usecase) ListAll() ([]models.Ticket, error) {
	return u.repo.ListAll()
}

func (u *Usecase) Create(ticket models.Ticket) (models.Ticket, error) {
	return u.repo.Create(ticket)
}

func (u *Usecase) FormByID(ticketID int) (models.Ticket, error) {
	return u.repo.FormByID(ticketID)
}

func (u *Usecase) ApproveByID(ticketID int) (models.Ticket, error) {
	return u.repo.ApproveByID(ticketID)
}

func (u *Usecase) RejectByID(ticketID int) (models.Ticket, error) {
	return u.repo.RejectByID(ticketID)
}

func (u *Usecase) DeleteByID(ticketID int) error {
	return u.repo.DeleteByID(ticketID)
}

func (u *Usecase) AddRoute(ticketID, routeID int) (models.Ticket, error) {
	return u.repo.AddRoute(ticketID, routeID)
}

func (u *Usecase) DeleteRoute(ticketID, routeID int) (models.Ticket, error) {
	return u.repo.AddRoute(ticketID, routeID)
}

func (u *Usecase) GetTicketDraftByCreatorID(creatorID int) *models.Ticket {
	return u.repo.GetTicketDraftByCreatorID(creatorID)
}
