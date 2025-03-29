package models

import (
	"github.com/gofrs/uuid/v5"
)

type Tag struct {
	ID   uuid.UUID `db:"tag_id" json:"tag_id"`
	Name string    `db:"name" json:"name"`
}
