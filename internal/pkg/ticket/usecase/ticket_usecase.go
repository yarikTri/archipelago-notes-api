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

func (u *Usecase) GetByID(ticketID uint32) (models.Ticket, error) {
	return u.repo.GetByID(ticketID)
}

func (u *Usecase) List() ([]models.Ticket, error) {
	return u.repo.List()
}

func (u *Usecase) Create(ticket models.Ticket) (models.Ticket, error) {
	return u.repo.Create(ticket)
}

func (u *Usecase) DeleteByID(ticketID uint32) error {
	return u.repo.DeleteByID(ticketID)
}
