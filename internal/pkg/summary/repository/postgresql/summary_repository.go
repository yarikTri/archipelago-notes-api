package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yarikTri/archipelago-notes-api/internal/common/repository"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

// PostgreSQL implements notes.Repository
type PostgreSQL struct {
	db *sqlx.DB
}

func NewPostgreSQL(db *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) UpdateSummaryTextRole(ID uuid.UUID, textWithRole, role string) error {
	query := fmt.Sprint(
		`UPDATE summ
		SET text_with_role = $2,
			role = $3
		WHERE id = $1;`, // TODO: check multiple update
	)
	if _, err := p.db.Exec(query, ID, textWithRole, role); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) SaveSummaryText(ID uuid.UUID, text string, active bool, detalization models.Detalization, platform string) (*models.Summary, error) {
	query := fmt.Sprint(
		`INSERT INTO summ (id, text, active, platform, detalization)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id)
		DO UPDATE SET text = $2, text_with_role = '', role = ''`, // TODO: check multiple update
	)
	if _, err := p.db.Exec(query, ID, text, active, platform, detalization); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &models.Summary{ID: ID, Text: text, Active: active}, nil
}

func (p *PostgreSQL) FinishSummary(ID uuid.UUID) error {
	query := fmt.Sprint(
		`UPDATE summ
		SET active = false
		WHERE id = $1;`,
	)
	if _, err := p.db.Exec(query, ID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetSummary(ID uuid.UUID) (*models.Summary, error) {
	query := fmt.Sprint(
		`SELECT id, text, active, text_with_role, role, platform, started_at, detalization
			FROM summ
			WHERE id = $1`,
	)

	var summary models.Summary
	if err := p.db.Get(&summary, query, ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %v", &repository.NotFoundError{ID: ID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &summary, nil
}
