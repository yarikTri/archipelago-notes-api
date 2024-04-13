package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (p *PostgreSQL) GetByID(noteID uuid.UUID) (*models.Note, error) {
	query := fmt.Sprint(
		`SELECT id, dir_id, automerge_url, title
			FROM note
			WHERE id = $1`,
	)

	var note models.Note
	if err := p.db.Get(&note, query, noteID.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %v", &repository.NotFoundError{ID: noteID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &note, nil
}

func (p *PostgreSQL) List() ([]*models.Note, error) {
	query := fmt.Sprint(
		`SELECT id, dir_id, automerge_url, title 
			FROM note`,
	)

	var notes []*models.Note
	if err := p.db.Select(&notes, query); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return notes, nil
}

func (p *PostgreSQL) ListByDirIds(dirIDs []int) ([]*models.Note, error) {
	query := fmt.Sprint(
		`SELECT id, automerge_url, title 
			FROM note
			WHERE dir_id = ANY($1)`,
	)

	var notes []*models.Note
	if err := p.db.Select(&notes, query, pq.Array(dirIDs)); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return notes, nil
}

func (p *PostgreSQL) Create(dirID int, automergeUrl, title string) (*models.Note, error) {
	query := fmt.Sprint(
		`INSERT INTO note (dir_id, automerge_url, title)
			VALUES ($1, $2, $3)
			RETURNING id`,
	)

	var noteID string
	row := p.db.QueryRow(query, dirID, automergeUrl, title)
	if err := row.Scan(&noteID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	uuid, _ := uuid.FromString(noteID)

	return &models.Note{ID: uuid, DirID: dirID, AutomergeURL: automergeUrl, Title: title}, nil
}

func (p *PostgreSQL) Update(note models.Note) (*models.Note, error) {
	query := fmt.Sprint(
		`UPDATE note
			SET dir_id = $1 automerge_url = $2, title = $3
			WHERE id = $4
			RETURNING id`,
	)

	_, err := p.db.Exec(query, note.DirID, note.AutomergeURL, note.Title, note.ID.String())
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &note, nil
}

func (p *PostgreSQL) DeleteByID(noteID uuid.UUID) error {
	query := fmt.Sprint(
		`DELETE
		FROM note
		WHERE id = $1`,
	)

	resExec, err := p.db.Exec(query, noteID.String())
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &repository.NotFoundError{ID: noteID})
	}

	return nil
}
