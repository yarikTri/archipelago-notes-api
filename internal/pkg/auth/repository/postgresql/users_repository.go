package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yarikTri/archipelago-notes-api/internal/common/repository"
)

// UsersRepository implements auth.UsersRepository
type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

type userIDAndPasswordRaw struct {
	ID       string `db:"id"`
	Password string `db:"password_hash"`
}

func (ur *UsersRepository) GetUserIDAndPasswordByEmail(email string) (uuid.UUID, string, error) {
	query := fmt.Sprint(
		`SELECT id, password_hash
			FROM "user"
			WHERE email = $1`,
	)

	var creds userIDAndPasswordRaw
	if err := ur.db.Get(&creds, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Max, "", fmt.Errorf("(repo) user not found: %v", err)
		}

		return uuid.Max, "", fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	id, _ := uuid.FromString(creds.ID)

	return id, creds.Password, nil
}

func (ur *UsersRepository) CreateUser(email, name, passwordHash string) (uuid.UUID, error) {
	query := fmt.Sprint(
		`INSERT INTO "user" (email, name, password_hash) VALUES ($1, $2, $3) RETURNING id`,
	)

	var userID string
	if err := ur.db.QueryRow(query, email, name, passwordHash).Scan(&userID); err != nil {
		return uuid.Max, fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	return uuid.FromString(userID)
}

func (ur *UsersRepository) DeleteUser(userID uuid.UUID) error {
	query := fmt.Sprint(
		`DELETE FROM user WHERE id = $1`,
	)

	resExec, err := ur.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	deleted, err := resExec.RowsAffected()
	if err != nil {
		return fmt.Errorf("(repo) failed to check RowsAffected: %w", err)
	}

	if deleted == 0 {
		return fmt.Errorf("(repo): %w", &repository.NotFoundError{ID: userID})
	}

	return nil
}
