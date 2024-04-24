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

// PostgreSQL implements users.Repository
type PostgreSQL struct {
	db *sqlx.DB
}

func NewPostgreSQL(db *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) GetByID(userID uuid.UUID) (*models.User, error) {
	query := fmt.Sprint(
		`SELECT u.id, u.email, u.name, urd.root_dir_id as root_dir_id
			FROM "user" u
				INNER JOIN user_root_dir urd ON u.id = urd.user_id
			WHERE u.id = $1`,
	)

	var user models.User
	if err := p.db.Get(&user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %v", &repository.NotFoundError{ID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &user, nil
}

func (p *PostgreSQL) Search(searchQuery string) ([]*models.User, error) {
	query := fmt.Sprint(
		`SELECT
				u.id as id,
				u.email as email,
				u.name as name,
				urd.root_dir_id as root_dir_id
			FROM "user" u
				INNER JOIN user_root_dir urd ON u.id = urd.user_id
			WHERE lower(u.email) LIKE '%' || lower($1) || '%' OR
			      LOWER(u.name) LIKE '%' || lower($1) || '%'`,
	)

	var users []*models.User
	if err := p.db.Select(&users, query, searchQuery); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return users, nil
}

func (p *PostgreSQL) SetRootDirByID(userID uuid.UUID, rootID int) error {
	query := fmt.Sprint(
		`INSERT INTO user_root_dir (user_id, root_dir_id) VALUES ($1, $2)
			ON CONFLICT (user_id)
			DO UPDATE SET root_dir_id = $2`,
	)

	_, err := p.db.Exec(query, userID.String(), rootID)
	return err
}
