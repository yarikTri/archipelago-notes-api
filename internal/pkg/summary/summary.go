package summary

import (
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type Usecase interface {
	SaveSummaryText(ID uuid.UUID, text string, active bool, detalization models.Detalization, platform string) (*models.Summary, error)
	UpdateSummaryTextRole(ID uuid.UUID, textWithRole, role string) error
	FinishSummary(ID uuid.UUID) error
	GetSummary(ID uuid.UUID) (*models.Summary, error)
	GetActiveSummaries() ([]models.Summary, error)
	UpdateName(ID uuid.UUID, name string) error
}

type Repository interface {
	SaveSummaryText(ID uuid.UUID, text string, active bool, detalization models.Detalization, platform string) (*models.Summary, error)
	UpdateSummaryTextRole(ID uuid.UUID, textWithRole, role string) error
	FinishSummary(ID uuid.UUID) error
	GetSummary(ID uuid.UUID) (*models.Summary, error)
	GetActiveSummaries() ([]models.Summary, error)
	UpdateName(ID uuid.UUID, name string) error
}
