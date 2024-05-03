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
		`SELECT id, dir_id, automerge_url, title, creator_id, default_access
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

func (p *PostgreSQL) List(userID uuid.UUID) ([]*models.Note, error) {
	query := fmt.Sprint(
		`SELECT
				id as id,
				dir_id as dir_id,
				automerge_url as automerge_url,
				title as title,
				creator_id as creator_id,
				default_access as default_access,
				'ma' as access
			FROM note
			WHERE creator_id = $1
			UNION ALL
			SELECT
			    n.id as id,
			    n.dir_id as dir_id,
			    n.automerge_url as automerge_url,
			    n.title as title,
			    n.creator_id as creator_id,
			    n.default_access as default_access,
			    na.access as access
			FROM note n INNER JOIN note_access na ON n.id = na.note_id
			WHERE na.user_id = $1 AND na.access <> 'e'`,
	)

	var notes []*models.Note
	if err := p.db.Select(&notes, query, userID.String()); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return notes, nil
}

func (p *PostgreSQL) ListByDirIds(dirIDs []int) ([]*models.Note, error) {
	query := fmt.Sprint(
		`SELECT id, dir_id, automerge_url, title, creator_id, default_access
			FROM note
			WHERE dir_id = ANY($1)`,
	)

	var notes []*models.Note
	if err := p.db.Select(&notes, query, pq.Array(dirIDs)); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return notes, nil
}

func (p *PostgreSQL) Create(dirID int, automergeUrl, title string, creatorID uuid.UUID) (*models.Note, error) {
	query := fmt.Sprint(
		`INSERT INTO note (dir_id, automerge_url, title, creator_id, default_access)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`,
	)

	defaultDefaultAccess := models.EmptyNoteAccess

	var noteID string
	row := p.db.QueryRow(query, dirID, automergeUrl, title, creatorID.String(), defaultDefaultAccess.String())
	if err := row.Scan(&noteID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	_uuid, _ := uuid.FromString(noteID)

	return &models.Note{
		ID:            _uuid,
		DirID:         dirID,
		AutomergeURL:  automergeUrl,
		Title:         title,
		CreatorID:     creatorID,
		DefaultAccess: defaultDefaultAccess.String(),
	}, nil
}

func (p *PostgreSQL) Update(note models.Note) (*models.Note, error) {
	query := fmt.Sprint(
		`UPDATE note
			SET dir_id = $1, automerge_url = $2, title = $3, default_access = $4
			WHERE id = $5`,
	)

	if _, err := p.db.Exec(query, note.DirID, note.AutomergeURL, note.Title, note.DefaultAccess, note.ID.String()); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return p.GetByID(note.ID)
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

type noteAccessInfo struct {
	CreatorID     string  `db:"creator_id"`
	DefaultAccess string  `db:"default_access"`
	Access        *string `db:"access"`
}

func (p *PostgreSQL) GetUserAccess(noteID uuid.UUID, userID uuid.UUID) (models.NoteAccess, error) {
	query := fmt.Sprint(
		`SELECT
			n.creator_id AS creator_id,
			n.default_access AS default_access,
			a.access AS access
		FROM note n LEFT JOIN (
			SELECT note_id, access
			FROM note_access
			WHERE note_id = $1 AND user_id = $2
		) a ON n.id = a.note_id
		WHERE n.id = $1`,
	)

	var access noteAccessInfo
	if err := p.db.Get(&access, query, noteID.String(), userID.String()); err != nil {
		return models.UndefinedNoteAccess, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	if access.CreatorID == userID.String() {
		return models.ManageAccessNoteAccess, nil
	}

	if access.Access != nil {
		return models.NoteAccessFromString(*access.Access), nil
	}

	return models.NoteAccessFromString(access.DefaultAccess), nil
}

func (p *PostgreSQL) SetUserAccess(noteID uuid.UUID, userID uuid.UUID, access models.NoteAccess) error {
	query := fmt.Sprint(
		`INSERT INTO note_access
		(note_id, user_id, access)
		VALUES ($1, $2, $3)`,
	)

	row := p.db.QueryRow(query, noteID.String(), userID.String(), access.String())
	if err := row.Scan(&access); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) AttachNoteToSummary(summID, noteID uuid.UUID) error {
	query := fmt.Sprint(
		`INSERT INTO summ_to_note (summ_id, note_id)
		VALUES ($1, $2);`,
	)
	if _, err := p.db.Exec(query, summID, noteID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DettachNoteFromSummary(summID, noteID uuid.UUID) error {
	query := fmt.Sprint(
		`delete from summ_to_note where summ_id = $1 AND note_id = $2`,
	)
	if _, err := p.db.Exec(query, summID, noteID); err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetSummaryListByNote(noteID uuid.UUID) ([]models.SummaryIDStatus, error) {
	query := fmt.Sprint(
		`SELECT summ.id, summ.active
		FROM summ_to_note
			INNER JOIN summ ON summ_to_note.summ_id = summ.id
		WHERE summ_to_note.note_id = $1;`,
	)

	var notesFromQuery []*models.SummaryIDStatus
	if err := p.db.Select(&notesFromQuery, query, noteID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	notes := make([]models.SummaryIDStatus, len(notesFromQuery))
	for i, notePtr := range notesFromQuery {
		notes[i] = *notePtr
	}

	return notes, nil
}
