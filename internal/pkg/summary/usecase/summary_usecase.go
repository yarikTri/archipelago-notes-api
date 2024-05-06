package usecase

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/summary"
)

// Usecase implements notes.Usecase
type Usecase struct {
	repo summary.Repository
}

func NewUsecase(rr summary.Repository) *Usecase {
	return &Usecase{
		repo: rr,
	}
}

func (u *Usecase) SaveSummaryText(ID uuid.UUID, text string, active bool, detalization models.Detalization, platform string) (*models.Summary, error) {
	return u.repo.SaveSummaryText(ID, text, active, detalization, platform)
}

func (u *Usecase) UpdateSummaryTextRole(ID uuid.UUID, textWithRole, role string) error {
	return u.repo.UpdateSummaryTextRole(ID, textWithRole, role)
}

func (u *Usecase) GetSummary(ID uuid.UUID) (*models.Summary, error) {
	return u.repo.GetSummary(ID)
}

func (u *Usecase) GetActiveSummaries() ([]models.Summary, error) {
	return u.repo.GetActiveSummaries()
}

func (u *Usecase) FinishSummary(ID uuid.UUID) error {
	return u.repo.FinishSummary(ID)
}

func (u *Usecase) UpdateName(ID uuid.UUID, name string) error {
	return u.repo.UpdateName(ID, name)
}
