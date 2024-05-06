package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Detalization uint8

const (
	Default Detalization = iota
	Short
	Long
)

func DetalizationFromString(d string) Detalization {
	switch d {
	case "":
	case "Средняя":
		return Default
	case "Краткая":
		return Short
	case "Развернутая":
		return Long
	}

	return Default
}

func (d Detalization) String() string {
	switch d {
	case Default:
		return "Средняя"
	case Short:
		return "Краткая"
	case Long:
		return "Развернутая"
	}

	return "Средняя"
}

type Summary struct {
	ID           uuid.UUID    `db:"id"`
	Text         string       `db:"text"`
	TextWithRole string       `db:"text_with_role"`
	Role         string       `db:"role"`
	Active       bool         `db:"active"`
	Platform     string       `db:"platform"`
	StartedAt    time.Time    `db:"started_at"`
	Detalization Detalization `db:"detalization"`
	Name         string       `db:"name"`
}

type SummaryIDStatus struct {
	ID     uuid.UUID `db:"id"`
	Active bool      `db:"active"`
}

func (s *Summary) ToTransfer() *SummaryTransfer {
	return &SummaryTransfer{
		ID:           s.ID.String(),
		Text:         s.Text,
		Active:       s.Active,
		TextWithRole: s.TextWithRole,
		Role:         s.Role,
		Platform:     s.Platform,
		StartedAt:    s.StartedAt,
		Detalization: s.Detalization.String(),
		Name:         s.Name,
	}
}

type SummaryTransfer struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	TextWithRole string    `json:"text_with_role"`
	Active       bool      `json:"active"`
	Role         string    `json:"role"`
	Platform     string    `json:"platform"`
	StartedAt    time.Time `json:"started_at"`
	Detalization string    `json:"detalization"`
	Name         string    `json:"name"`
}
