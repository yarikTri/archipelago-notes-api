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

// PostgreSQL implements tag.TagRepository
type PostgreSQL struct {
	db *sqlx.DB
}

func NewPostgreSQL(db *sqlx.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

// Helper functions to reduce code duplication

type QueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

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

func (p *PostgreSQL) getTagByID(q QueryExecutor, id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := q.QueryRow(`
		SELECT tag_id, name, user_id
		FROM tag
		WHERE tag_id = $1`, id).Scan(&tag.ID, &tag.Name, &tag.UserID)

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
		SELECT tag_id, name, user_id
		FROM tag
		WHERE name = $1`, name).Scan(&tag.ID, &tag.Name, &tag.UserID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to check existing tag: %w", err)
	}
	return &tag, nil
}

func (p *PostgreSQL) getTagByNameAndUserID(tx *sql.Tx, name string, userID uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := tx.QueryRow(`
		SELECT tag_id, name, user_id
		FROM tag
		WHERE name = $1 and user_id = $2`, name, userID).Scan(&tag.ID, &tag.Name, &tag.UserID)

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

func (p *PostgreSQL) CreateAndLinkTag(name string, noteID, userID uuid.UUID) (*models.Tag, error) {
	var result *models.Tag
	err := p.withTransaction(func(tx *sql.Tx) error {
		// Check if note exists and get its creator ID
		var noteExists bool
		err := tx.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM note WHERE id = $1
			)`, noteID).Scan(&noteExists)
		if err != nil {
			return fmt.Errorf("(repo) failed to check note existence: %w", err)
		}
		if !noteExists {
			return &errors.NoteNotFoundError{ID: noteID}
		}

		// Check if tag with this name is already linked to this note
		var existingTagID uuid.UUID
		err = tx.QueryRow(`
			SELECT t.tag_id
			FROM tag t
			JOIN tag_to_note ttn ON t.tag_id = ttn.tag_id
			WHERE t.name = $1 AND ttn.note_id = $2
			LIMIT 1`, name, noteID).Scan(&existingTagID)
		if err == nil {
			// Tag already exists and is linked to this note
			return &errors.TagNameExistsError{
				Name: name,
				ID:   existingTagID,
			}
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("(repo) failed to check existing tag: %w", err)
		}

		// Create a new tag
		id, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("(repo) failed to generate uuid: %w", err)
		}

		// Create a new tag
		_, err = tx.Exec(`
			INSERT INTO tag (tag_id, user_id, name)
			VALUES ($1, $2, $3)`,
			id, userID, name)
		if err != nil {
			return fmt.Errorf("(repo) failed to create tag: %w", err)
		}

		// Link the tag to the note
		_, err = tx.Exec(`
			INSERT INTO tag_to_note (tag_id, note_id)
			VALUES ($1, $2)`,
			id, noteID)
		if err != nil {
			return fmt.Errorf("(repo) failed to link new tag: %w", err)
		}

		result = &models.Tag{
			ID:     id,
			Name:   name,
			UserID: userID,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgreSQL) LinkTagToNote(tagID uuid.UUID, noteID uuid.UUID) error {
	return p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag exists
		_, err := p.getTagByID(tx, tagID)
		if err != nil {
			return err
		}

		// Check if note exists
		var noteExists bool
		err = tx.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM note WHERE id = $1
			)`, noteID).Scan(&noteExists)
		if err != nil {
			return fmt.Errorf("(repo) failed to check note existence: %w", err)
		}
		if !noteExists {
			return &errors.NoteNotFoundError{ID: noteID}
		}

		// Check if tag is already linked to the note
		var exists bool
		err = tx.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM tag_to_note
				WHERE tag_id = $1 AND note_id = $2
			)`, tagID, noteID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("(repo) failed to check tag link existence: %w", err)
		}
		if exists {
			return &errors.TagLinkExistsError{}
		}

		// Link tag to note
		_, err = tx.Exec(`
			INSERT INTO tag_to_note (tag_id, note_id)
			VALUES ($1, $2)`,
			tagID, noteID)
		if err != nil {
			return fmt.Errorf("(repo) failed to link tag to note: %w", err)
		}

		return nil
	})
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

func (p *PostgreSQL) UpdateTag(ID uuid.UUID, name string, userID uuid.UUID) (*models.Tag, error) {
	var result *models.Tag
	err := p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag exists
		existingTag, err := p.getTagByID(tx, ID)
		if err != nil {
			return err
		}

		// If the new name is the same as the current name, no need to update
		if existingTag.Name == name {
			result = existingTag
			return nil
		}

		// Check if new name already exists
		conflictingTag, err := p.getTagByNameAndUserID(tx, name, userID)
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

		result = &models.Tag{
			ID:     ID,
			Name:   name,
			UserID: existingTag.UserID,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// func (p *PostgreSQL) UpdateTagForNote(tagID uuid.UUID, noteID uuid.UUID, newName string) (*models.Tag, error) {
// 	var result *models.Tag
// 	err := p.withTransaction(func(tx *sql.Tx) error {
// 		// Check if tag exists and is linked to the note
// 		var tag models.Tag
// 		err := tx.QueryRow(`
// 			SELECT t.tag_id, t.name
// 			FROM tag t
// 			JOIN tag_to_note ttn ON t.tag_id = ttn.tag_id
// 			WHERE t.tag_id = $1 AND ttn.note_id = $2`, tagID, noteID).Scan(&tag.ID, &tag.Name)
// 		if err == sql.ErrNoRows {
// 			return &errors.TagNotFoundError{ID: tagID}
// 		}
// 		if err != nil {
// 			return fmt.Errorf("(repo) failed to check tag existence: %w", err)
// 		}

// 		// Check if new name is already used by another tag
// 		var exists bool
// 		err = tx.QueryRow(`
// 			SELECT EXISTS (
// 				SELECT 1 FROM tag WHERE name = $1
// 			)`, newName).Scan(&exists)
// 		if err != nil {
// 			return fmt.Errorf("(repo) failed to check tag name existence: %w", err)
// 		}
// 		if exists {
// 			return &errors.TagNameExistsError{}
// 		}

// 		// Create new tag with the new name
// 		newTagID, err := uuid.NewV4()
// 		if err != nil {
// 			return fmt.Errorf("(repo) failed to generate uuid: %w", err)
// 		}

// 		_, err = tx.Exec(`
// 			INSERT INTO tag (tag_id, name)
// 			VALUES ($1, $2)`, newTagID, newName)
// 		if err != nil {
// 			return fmt.Errorf("(repo) failed to create new tag: %w", err)
// 		}

// 		// Link new tag to the note
// 		_, err = tx.Exec(`
// 			INSERT INTO tag_to_note (tag_id, note_id)
// 			VALUES ($1, $2)`, newTagID, noteID)
// 		if err != nil {
// 			return fmt.Errorf("(repo) failed to link new tag to note: %w", err)
// 		}

// 		// Remove old tag link from the note
// 		_, err = tx.Exec(`
// 			DELETE FROM tag_to_note
// 			WHERE tag_id = $1 AND note_id = $2`, tagID, noteID)
// 		if err != nil {
// 			return fmt.Errorf("(repo) failed to unlink old tag from note: %w", err)
// 		}

// 		// Check if old tag has any remaining relations
// 		count, err := p.getTagRelationsCount(tx, tagID)
// 		if err != nil {
// 			return err
// 		}

// 		// If no relations left, delete the old tag
// 		if count == 0 {
// 			_, err = tx.Exec(`
// 				DELETE FROM tag
// 				WHERE tag_id = $1`, tagID)
// 			if err != nil {
// 				return fmt.Errorf("(repo) failed to delete old tag: %w", err)
// 			}
// 		}

// 		result = &models.Tag{
// 			ID:   newTagID,
// 			Name: newName,
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

func (p *PostgreSQL) GetNotesByTag(tagID uuid.UUID) ([]models.Note, error) {
	// TODO: maybe check that all returned notes are visible by user.

	// Check if tag exists first
	var tag models.Tag
	err := p.db.Get(&tag, `
		SELECT tag_id, name, user_id
		FROM tag
		WHERE tag_id = $1`, tagID)

	if err == sql.ErrNoRows {
		return nil, &errors.TagNotFoundError{ID: tagID}
	}
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to check tag existence: %w", err)
	}

	var notes []models.Note
	err = p.db.Select(&notes, `
		SELECT n.id, n.dir_id, n.title, n.automerge_url, n.creator_id, n.default_access
		FROM note n
		JOIN tag_to_note ttn ON n.id = ttn.note_id
		WHERE ttn.tag_id = $1`, tagID)
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to get notes by tag: %w", err)
	}
	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

func (p *PostgreSQL) GetTagsByNoteForUser(noteID, userID uuid.UUID) ([]models.Tag, error) {
	// Check if note exists first
	var noteExists bool
	err := p.db.Get(&noteExists, `
		SELECT EXISTS (
			SELECT 1 FROM note WHERE id = $1
		)`, noteID)
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to check note existence: %w", err)
	}
	if !noteExists {
		return nil, &errors.NoteNotFoundError{ID: noteID}
	}

	// Get all tags linked to the note
	var tags []models.Tag
	err = p.db.Select(&tags, `
		SELECT t.tag_id, t.name, t.user_id
		FROM tag t
		JOIN tag_to_note ttn ON t.tag_id = ttn.tag_id
		WHERE ttn.note_id = $1 and t.user_id = $2
		ORDER BY t.name`, noteID, userID)

	if err != nil {
		return nil, fmt.Errorf("(repo) failed to get tags by note: %w", err)
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	return tags, nil
}

func (p *PostgreSQL) LinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return p.withTransaction(func(tx *sql.Tx) error {
		// Check if both tags exist
		_, err := p.getTagByID(tx, tag1ID)
		if err != nil {
			return &errors.TagNotFoundError{ID: tag1ID}
		}

		_, err = p.getTagByID(tx, tag2ID)
		if err != nil {
			return &errors.TagNotFoundError{ID: tag2ID}
		}

		// Check if tags are already linked
		var exists bool
		err = tx.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM tag_to_tag
				WHERE (tag_1_id = $1 AND tag_2_id = $2)
				OR (tag_1_id = $2 AND tag_2_id = $1)
			)`, tag1ID, tag2ID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("(repo) failed to check tag link existence: %w", err)
		}
		if exists {
			return &errors.TagLinkExistsError{}
		}

		// Create the link
		_, err = tx.Exec(`
			INSERT INTO tag_to_tag (tag_1_id, tag_2_id)
			VALUES ($1, $2)`, tag1ID, tag2ID)
		if err != nil {
			return fmt.Errorf("(repo) failed to create tag link: %w", err)
		}

		return nil
	})
}

func (p *PostgreSQL) UnlinkTags(tag1ID uuid.UUID, tag2ID uuid.UUID) error {
	return p.withTransaction(func(tx *sql.Tx) error {
		// Check if both tags exist
		_, err := p.getTagByID(tx, tag1ID)
		if err != nil {
			return &errors.TagNotFoundError{ID: tag1ID}
		}

		_, err = p.getTagByID(tx, tag2ID)
		if err != nil {
			return &errors.TagNotFoundError{ID: tag2ID}
		}

		// Check if tags are linked
		var exists bool
		err = tx.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM tag_to_tag
				WHERE (tag_1_id = $1 AND tag_2_id = $2)
				OR (tag_1_id = $2 AND tag_2_id = $1)
			)`, tag1ID, tag2ID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("(repo) failed to check tag link existence: %w", err)
		}
		if !exists {
			return &errors.TagLinkNotFoundError{}
		}

		// Remove the link
		_, err = tx.Exec(`
			DELETE FROM tag_to_tag
			WHERE (tag_1_id = $1 AND tag_2_id = $2)
			OR (tag_1_id = $2 AND tag_2_id = $1)`, tag1ID, tag2ID)
		if err != nil {
			return fmt.Errorf("(repo) failed to remove tag link: %w", err)
		}

		return nil
	})
}

func (p *PostgreSQL) GetLinkedTagsForUser(tagID, userID uuid.UUID) ([]models.LinkedTag, error) {
	if err := p.validateUserTagExists(tagID, userID); err != nil {
		return nil, err
	}

	var tags []models.LinkedTag
	if err := p.db.Select(&tags, `
		SELECT t.tag_id, t.name, t.user_id, ttt.name link_name
		FROM tag t
			JOIN tag_to_tag ttt ON (ttt.tag_1_id = t.tag_id OR ttt.tag_2_id = t.tag_id)
		WHERE (ttt.tag_1_id = $1 OR ttt.tag_2_id = $1)
			AND t.tag_id != $1
			AND t.user_id = $2`, tagID, userID); err != nil {
		return nil, fmt.Errorf("(repo) failed to get linked tags: %w", err)
	}
	if tags == nil {
		tags = []models.LinkedTag{}
	}

	return tags, nil
}

func (p *PostgreSQL) UpdateTagsLinkName(tag1ID, tag2ID, userID uuid.UUID, linkName string) error {
	if err := p.validateUserTagExists(tag1ID, userID); err != nil {
		return err
	}
	if err := p.validateUserTagExists(tag2ID, userID); err != nil {
		return err
	}

	if _, err := p.db.Exec(`
		UPDATE tag_to_tag
		SET name = $1
		WHERE tag_1_id IN ($2, $3) AND tag_2_id IN ($2, $3)`,
		linkName, tag1ID, tag2ID); err != nil {
		return fmt.Errorf("(repo) failed to update tags link name: %w", err)
	}

	return nil
}

func (p *PostgreSQL) validateUserTagExists(tagID, userID uuid.UUID) error {
	var tag models.Tag
	err := p.db.Get(&tag, `
		SELECT tag_id, name, user_id
		FROM tag
		WHERE tag_id = $1`, tagID)

	if err != nil {
		if err == sql.ErrNoRows {
			return &errors.TagNotFoundError{ID: tagID}
		}
		return fmt.Errorf("(repo) failed to check tag existence: %w", err)
	}

	if tag.UserID != userID {
		return &errors.TagNotFoundError{ID: tagID}
	}

	return nil
}

// DeleteTag deletes a tag and all its relations
func (p *PostgreSQL) DeleteTag(tagID uuid.UUID) error {
	err := p.withTransaction(func(tx *sql.Tx) error {
		// Check if tag exists
		_, err := p.getTagByID(tx, tagID)
		if err != nil {
			return err
		}

		// Delete all tag-to-note relations
		_, err = tx.Exec(`
			DELETE FROM tag_to_note
			WHERE tag_id = $1`, tagID)
		if err != nil {
			return fmt.Errorf("(repo) failed to delete tag-to-note relations: %w", err)
		}

		// Delete all tag-to-tag relations
		_, err = tx.Exec(`
			DELETE FROM tag_to_tag
			WHERE tag_1_id = $1 OR tag_2_id = $1`, tagID)
		if err != nil {
			return fmt.Errorf("(repo) failed to delete tag-to-tag relations: %w", err)
		}

		// Delete the tag itself
		_, err = tx.Exec(`			DELETE FROM tag
			WHERE tag_id = $1`, tagID)
		if err != nil {
			return fmt.Errorf("(repo) failed to delete tag: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSQL) GetTagByID(tagID uuid.UUID) (*models.Tag, error) {
	return p.getTagByID(p.db.DB, tagID)
}
