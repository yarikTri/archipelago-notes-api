package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/repository"
)

// PostgreSQL implements tag.Repository
type PostgreSQL struct {
	db *sqlx.DB
}

func NewPostgreSQL(db *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if tag with this name exists
	var existingTag models.Tag
	err = tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE name = $1`, name).Scan(&existingTag.ID, &existingTag.Name)

	if err == nil {
		// Tag exists, just link it with the note
		_, err = tx.Exec(`
			INSERT INTO tag_to_note (tag_id, note_id)
			VALUES ($1, $2)
			ON CONFLICT (tag_id, note_id) DO NOTHING`,
			existingTag.ID, noteID)
		if err != nil {
			return nil, fmt.Errorf("(repo) failed to link existing tag: %w", err)
		}
	} else if err == sql.ErrNoRows {
		// Tag doesn't exist, create it and link
		id, err := uuid.NewV4()
		if err != nil {
			return nil, fmt.Errorf("(repo) failed to generate uuid: %w", err)
		}

		_, err = tx.Exec(`
			INSERT INTO tag (tag_id, name)
			VALUES ($1, $2)`,
			id, name)
		if err != nil {
			return nil, fmt.Errorf("(repo) failed to create tag: %w", err)
		}

		_, err = tx.Exec(`
			INSERT INTO tag_to_note (tag_id, note_id)
			VALUES ($1, $2)`,
			id, noteID)
		if err != nil {
			return nil, fmt.Errorf("(repo) failed to link new tag: %w", err)
		}

		existingTag = models.Tag{
			ID:   id,
			Name: name,
		}
	} else {
		return nil, fmt.Errorf("(repo) failed to check existing tag: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("(repo) failed to commit transaction: %w", err)
	}

	return &existingTag, nil
}

func (p *PostgreSQL) UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Remove the tag-note relation
	_, err = tx.Exec(`
		DELETE FROM tag_to_note
		WHERE tag_id = $1 AND note_id = $2`,
		tagID, noteID)
	if err != nil {
		return fmt.Errorf("(repo) failed to unlink tag from note: %w", err)
	}

	// Check if there are any remaining relations for this tag
	var count int
	err = tx.QueryRow(`
		SELECT COUNT(*)
		FROM tag_to_note
		WHERE tag_id = $1`,
		tagID).Scan(&count)
	if err != nil {
		return fmt.Errorf("(repo) failed to count tag relations: %w", err)
	}

	// If no relations left, delete the tag
	if count == 0 {
		_, err = tx.Exec(`
			DELETE FROM tag
			WHERE tag_id = $1`,
			tagID)
		if err != nil {
			return fmt.Errorf("(repo) failed to delete tag: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("(repo) failed to commit transaction: %w", err)
	}

	return nil
}

// Placeholder implementations for other interface methods
func (p *PostgreSQL) GetTag(ID uuid.UUID) (*models.Tag, error) {
	return nil, nil
}

func (p *PostgreSQL) GetAllTags() ([]models.Tag, error) {
	return nil, nil
}

func (p *PostgreSQL) UpdateTag(ID uuid.UUID, name string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if tag exists
	var existingTag models.Tag
	err = tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE tag_id = $1`, ID).Scan(&existingTag.ID, &existingTag.Name)

	if err == sql.ErrNoRows {
		return fmt.Errorf("(repo) tag not found: %w", &repository.NotFoundError{ID: ID})
	}
	if err != nil {
		return fmt.Errorf("(repo) failed to check existing tag: %w", err)
	}

	// If the new name is the same as the current name, no need to update
	if existingTag.Name == name {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("(repo) failed to commit transaction: %w", err)
		}
		return nil
	}

	// Check if new name already exists
	var conflictingTag models.Tag
	err = tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE name = $1`, name).Scan(&conflictingTag.ID, &conflictingTag.Name)

	if err == nil {
		return fmt.Errorf("(repo) tag with name '%s' already exists (ID: %s)", name, conflictingTag.ID)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("(repo) failed to check name uniqueness: %w", err)
	}

	// Update tag name
	_, err = tx.Exec(`
		UPDATE tag
		SET name = $1
		WHERE tag_id = $2`, name, ID)
	if err != nil {
		return fmt.Errorf("(repo) failed to update tag: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("(repo) failed to commit transaction: %w", err)
	}

	return nil
}

func (p *PostgreSQL) DeleteTag(ID uuid.UUID) error {
	return nil
}

func (p *PostgreSQL) UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if tag exists
	var existingTag models.Tag
	err = tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE tag_id = $1`, tagID).Scan(&existingTag.ID, &existingTag.Name)

	if err == sql.ErrNoRows {
		return fmt.Errorf("(repo) tag not found: %w", &repository.NotFoundError{ID: tagID})
	}
	if err != nil {
		return fmt.Errorf("(repo) failed to check existing tag: %w", err)
	}

	// If the new name is the same as the current name, no need to update
	if existingTag.Name == newName {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("(repo) failed to commit transaction: %w", err)
		}
		return nil
	}

	// Check if new name already exists
	var conflictingTag models.Tag
	err = tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE name = $1`, newName).Scan(&conflictingTag.ID, &conflictingTag.Name)

	if err == nil {
		// Tag with new name exists, update the link
		_, err = tx.Exec(`
			UPDATE tag_to_note
			SET tag_id = $1
			WHERE tag_id = $2 AND note_id = $3`,
			conflictingTag.ID, tagID, noteID)
		if err != nil {
			return fmt.Errorf("(repo) failed to update tag link: %w", err)
		}
	} else if err == sql.ErrNoRows {
		// Create new tag with the new name
		newTagID, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("(repo) failed to generate uuid: %w", err)
		}

		_, err = tx.Exec(`
			INSERT INTO tag (tag_id, name)
			VALUES ($1, $2)`,
			newTagID, newName)
		if err != nil {
			return fmt.Errorf("(repo) failed to create new tag: %w", err)
		}

		// Update the link to point to the new tag
		_, err = tx.Exec(`
			UPDATE tag_to_note
			SET tag_id = $1
			WHERE tag_id = $2 AND note_id = $3`,
			newTagID, tagID, noteID)
		if err != nil {
			return fmt.Errorf("(repo) failed to update tag link: %w", err)
		}
	} else {
		return fmt.Errorf("(repo) failed to check name uniqueness: %w", err)
	}

	// Check if old tag has any remaining relations
	var count int
	err = tx.QueryRow(`
		SELECT COUNT(*)
		FROM tag_to_note
		WHERE tag_id = $1`,
		tagID).Scan(&count)
	if err != nil {
		return fmt.Errorf("(repo) failed to count tag relations: %w", err)
	}

	// If no relations left, delete the old tag
	if count == 0 {
		_, err = tx.Exec(`
			DELETE FROM tag
			WHERE tag_id = $1`,
			tagID)
		if err != nil {
			return fmt.Errorf("(repo) failed to delete old tag: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("(repo) failed to commit transaction: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetNotesByTag(tagID uuid.UUID) ([]models.Note, error) {
	// Check if tag exists first
	var tag models.Tag
	err := p.db.Get(&tag, `
		SELECT tag_id, name
		FROM tag
		WHERE tag_id = $1`, tagID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("(repo) tag not found: %w", &repository.NotFoundError{ID: tagID})
	}
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to check tag existence: %w", err)
	}

	// Get all notes linked to the tag
	var notes []models.Note
	err = p.db.Select(&notes, `
		SELECT n.id, n.dir_id, n.title, n.automerge_url, n.creator_id, n.default_access
		FROM note n
		JOIN tag_to_note ttn ON n.id = ttn.note_id
		WHERE ttn.tag_id = $1`, tagID)

	if err != nil {
		return nil, fmt.Errorf("(repo) failed to get notes by tag: %w", err)
	}

	return notes, nil
}

func (p *PostgreSQL) GetTagsByNote(noteID uuid.UUID) ([]models.Tag, error) {
	// Get all tags linked to the note
	var tags []models.Tag
	err := p.db.Select(&tags, `
		SELECT t.tag_id, t.name
		FROM tag t
		JOIN tag_to_note ttn ON t.tag_id = ttn.tag_id
		WHERE ttn.note_id = $1
		ORDER BY t.name`, noteID)

	if err != nil {
		return nil, fmt.Errorf("(repo) failed to get tags by note: %w", err)
	}

	return tags, nil
}
