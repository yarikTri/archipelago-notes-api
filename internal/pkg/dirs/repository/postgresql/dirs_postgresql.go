package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/yarikTri/archipelago-notes-api/internal/common/repository"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

// PostgreSQL implements dirs.Repository
type PostgreSQL struct {
	db *sqlx.DB
}

func NewPostgreSQL(db *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) GetByID(dirID int) (*models.Dir, error) {
	query := fmt.Sprint(
		`SELECT id, name, SUBPATH(path, 0, -1) as subpath
			FROM dir
			WHERE id = $1`,
	)

	var dir models.Dir
	if err := p.db.Get(&dir, query, dirID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %v", &repository.NotFoundError{ID: dirID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &dir, nil
}

func (p *PostgreSQL) GetSubTreeDirsByID(dirID int) ([]*models.Dir, error) {
	query := fmt.Sprint(
		`SELECT id, name, SUBPATH(path, 0, -1) as subpath
			FROM dir
			WHERE path <@ (SELECT path FROM dir WHERE id = $1)`,
	)

	var dirs []*models.Dir
	if err := p.db.Select(&dirs, query, dirID); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return dirs, nil
}

func (p *PostgreSQL) Create(parentDirID int, name string) (*models.Dir, error) {
	var query string

	var id int
	var createdName string
	var path string
	var row *sql.Row
	if parentDirID == 0 {
		query = fmt.Sprint(
			`INSERT INTO dir (name)
			VALUES ($1)
			RETURNING id, name, SUBPATH(path, 0, -1) as subpath`,
		)
		row = p.db.QueryRow(query, name)
	} else {
		query = fmt.Sprint(
			`INSERT INTO dir (name, path)
			VALUES ($1, (SELECT path FROM dir WHERE id = $2))
			RETURNING id, name, SUBPATH(path, 0, -1) as subpath`,
		)
		row = p.db.QueryRow(query, name, parentDirID)
	}

	if err := row.Scan(&id, &createdName, &path); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &models.Dir{ID: id, Name: createdName, Path: path}, nil
}

func (p *PostgreSQL) Update(dir *models.Dir) (*models.Dir, error) {
	query := fmt.Sprint(
		`UPDATE dir
			SET name = $1, path = $2
			WHERE id = $3
			RETURNING id, name, SUBPATH(path, 0, -1) as subpath`,
	)

	var id int
	var name string
	var path string
	row := p.db.QueryRow(query, dir.Name, dir.Path, dir.ID)
	if err := row.Scan(&id, &name, &path); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &models.Dir{ID: id, Name: name, Path: path}, nil
}

func (p *PostgreSQL) DeleteByID(dirID int) error {
	query := fmt.Sprint(
		`DELETE
		FROM dir
		WHERE id = $1`,
	)

	resExec, err := p.db.Exec(query, dirID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &repository.NotFoundError{ID: dirID})
	}

	return nil
}
