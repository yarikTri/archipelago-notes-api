package repository

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
)

type NotFoundError struct {
	ID uuid.UUID
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Row with id %s not found", e.ID.String())
}
