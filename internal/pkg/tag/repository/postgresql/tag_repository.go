package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/errors"
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

// Helper functions to reduce code duplication

// withTransaction executes a function within a transaction and handles commit/rollback
func (p *PostgreSQL) withTransaction(fn func(*sql.Tx) error) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("(repo) failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("(repo) failed to commit transaction: %w", err)
	}
	return nil
}

// getTagByID retrieves a tag by its ID
func (p *PostgreSQL) getTagByID(tx *sql.Tx, id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE tag_id = $1`, id).Scan(&tag.ID, &tag.Name)

	if err == sql.ErrNoRows {
		return nil, &errors.TagNotFoundError{ID: id}
	}
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to check existing tag: %w", err)
	}
	return &tag, nil
}

// getTagByName retrieves a tag by its name
func (p *PostgreSQL) getTagByName(tx *sql.Tx, name string) (*models.Tag, error) {
	var tag models.Tag
	err := tx.QueryRow(`
		SELECT tag_id, name
		FROM tag
		WHERE name = $1`, name).Scan(&tag.ID, &tag.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to check existing tag: %w", err)
	}
	return &tag, nil
}

// getTagRelationsCount counts how many notes are linked to a tag
func (p *PostgreSQL) getTagRelationsCount(tx *sql.Tx, tagID uuid.UUID) (int, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(*)
		FROM tag_to_note
		WHERE tag_id = $1`,
		tagID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("(repo) failed to count tag relations: %w", err)
	}
	return count, nil
}

// Main repository methods

func (p *PostgreSQL) CreateAndLinkTag(name string, noteID uuid.UUID) (*models.Tag, error) {
	var result *models.Tag
	err := p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag with this name exists
		existingTag, err := p.getTagByName(tx, name)
		if err != nil {
			return err
		}

		if existingTag != nil {
			// Tag exists, just link it with the note
			_, err = tx.Exec(`
				INSERT INTO tag_to_note (tag_id, note_id)
				VALUES ($1, $2)
				ON CONFLICT (tag_id, note_id) DO NOTHING`,
				existingTag.ID, noteID)
			if err != nil {
				return fmt.Errorf("(repo) failed to link existing tag: %w", err)
			}
			result = existingTag
		} else {
			// Tag doesn't exist, create it and link
			id, err := uuid.NewV4()
			if err != nil {
				return fmt.Errorf("(repo) failed to generate uuid: %w", err)
			}

			_, err = tx.Exec(`
				INSERT INTO tag (tag_id, name)
				VALUES ($1, $2)`,
				id, name)
			if err != nil {
				return fmt.Errorf("(repo) failed to create tag: %w", err)
			}

			_, err = tx.Exec(`
				INSERT INTO tag_to_note (tag_id, note_id)
				VALUES ($1, $2)`,
				id, noteID)
			if err != nil {
				return fmt.Errorf("(repo) failed to link new tag: %w", err)
			}

			result = &models.Tag{
				ID:   id,
				Name: name,
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgreSQL) UnlinkTagFromNote(tagID uuid.UUID, noteID uuid.UUID) error {
	return p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag exists first
		_, err := p.getTagByID(tx, tagID)
		if err != nil {
			return err
		}

		// Remove the tag-note relation
		result, err := tx.Exec(`
			DELETE FROM tag_to_note
			WHERE tag_id = $1 AND note_id = $2`,
			tagID, noteID)
		if err != nil {
			return fmt.Errorf("(repo) failed to unlink tag from note: %w", err)
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("(repo) failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			return &errors.TagLinkNotFoundError{
				TagID:  tagID,
				NoteID: noteID,
			}
		}

		// Check if there are any remaining relations for this tag
		count, err := p.getTagRelationsCount(tx, tagID)
		if err != nil {
			return err
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

		return nil
	})
}

func (p *PostgreSQL) UpdateTag(ID uuid.UUID, name string) error {
	return p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag exists
		existingTag, err := p.getTagByID(tx, ID)
		if err != nil {
			return err
		}

		// If the new name is the same as the current name, no need to update
		if existingTag.Name == name {
			return nil
		}

		// Check if new name already exists
		conflictingTag, err := p.getTagByName(tx, name)
		if err != nil {
			return err
		}
		if conflictingTag != nil {
			return &errors.TagNameExistsError{
				Name: name,
				ID:   conflictingTag.ID,
			}
		}

		// Update tag name
		_, err = tx.Exec(`
			UPDATE tag
			SET name = $1
			WHERE tag_id = $2`, name, ID)
		if err != nil {
			return fmt.Errorf("(repo) failed to update tag: %w", err)
		}

		return nil
	})
}

func (p *PostgreSQL) UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) error {
	return p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag exists
		existingTag, err := p.getTagByID(tx, tagID)
		if err != nil {
			return err
		}

		// If the new name is the same as the current name, no need to update
		if existingTag.Name == newName {
			return nil
		}

		// Check if new name already exists
		conflictingTag, err := p.getTagByName(tx, newName)
		if err != nil {
			return err
		}

		if conflictingTag != nil {
			// Tag with new name exists, update the link
			_, err = tx.Exec(`
				UPDATE tag_to_note
				SET tag_id = $1
				WHERE tag_id = $2 AND note_id = $3`,
				conflictingTag.ID, tagID, noteID)
			if err != nil {
				return fmt.Errorf("(repo) failed to update tag link: %w", err)
			}
		} else {
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
		}

		// Check if old tag has any remaining relations
		count, err := p.getTagRelationsCount(tx, tagID)
		if err != nil {
			return err
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

		return nil
	})
}

func (p *PostgreSQL) GetNotesByTag(tagID uuid.UUID) ([]models.Note, error) {
	// Check if tag exists first
	var tag models.Tag
	err := p.db.Get(&tag, `
		SELECT tag_id, name
		FROM tag
		WHERE tag_id = $1`, tagID)

	if err == sql.ErrNoRows {
		return nil, &errors.TagNotFoundError{ID: tagID}
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
