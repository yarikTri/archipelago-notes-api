package usecase

import (
	"time"

	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"
)

// Usecase implements ticket.Usecase
type Usecase struct {
	repo      ticket.Repository
	draftRepo ticket.DraftRepository
	routeRepo route.Repository
}

func NewUsecase(tr ticket.Repository, tdr ticket.DraftRepository, rr route.Repository) *Usecase {
	return &Usecase{
		repo:      tr,
		draftRepo: tdr,
		routeRepo: rr,
	}
}

func (u *Usecase) GetByID(ticketID int) (models.Ticket, error) {
	return u.repo.GetByID(ticketID)
}

func (u *Usecase) List() ([]models.Ticket, error) {
	return u.repo.List()
}

func (u *Usecase) Create(ticket models.Ticket) (models.Ticket, error) {
	return u.draftRepo.SetTicketDraft(ticket, true)
}

func (u *Usecase) FormDraft(creatorID int) (models.Ticket, error) {
	ticket, err := u.draftRepo.GetTicketDraft(creatorID)
	if err != nil {
		return models.Ticket{}, err
	}

	formedTicket, err := u.repo.Create(
		models.Ticket{
			State:     models.FORMED_STATE,
			FormTime:  time.Now(),
			CreatorID: ticket.CreatorID,
			Routes:    ticket.Routes,
		},
	)
	if err != nil {
		return models.Ticket{}, err
	}

	if err := u.draftRepo.DelTicketDraft(formedTicket.CreatorID); err != nil {
		return models.Ticket{}, err
	}

	return formedTicket, err
}

func (u *Usecase) ApproveByID(ticketID, moderatorID int) (models.Ticket, error) {
	return u.repo.ApproveByID(ticketID, moderatorID)
}

func (u *Usecase) RejectByID(ticketID, moderatorID int) (models.Ticket, error) {
	return u.repo.RejectByID(ticketID, moderatorID)
}

func (u *Usecase) DeleteByID(ticketID int) error {
	return u.repo.DeleteByID(ticketID)
}

func (u *Usecase) GetDraft(creatorID int) (models.Ticket, error) {
	return u.draftRepo.GetTicketDraft(creatorID)
}

func (u *Usecase) DeleteDraft(creatorID int) error {
	return u.draftRepo.DelTicketDraft(creatorID)
}

func (u *Usecase) AddRoute(creatorID, routeID int) (models.Ticket, error) {
	route, err := u.routeRepo.GetByID(routeID)
	if err != nil {
		return models.Ticket{}, err
	}

	return u.draftRepo.AddRoute(creatorID, route)
}

func (u *Usecase) DeleteRoute(creatorID, routeID int) (models.Ticket, error) {
	return u.draftRepo.DeleteRoute(creatorID, routeID)
}

func (u *Usecase) FinalizeWriting(ticketID int) (models.Ticket, error) {
	return u.repo.FinalizeWriting(ticketID)
}
