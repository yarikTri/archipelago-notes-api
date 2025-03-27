package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Tag struct {
	ID        uuid.UUID `db:"tag_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}
