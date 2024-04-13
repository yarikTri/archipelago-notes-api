package repository

import (
	"fmt"
)

type NotFoundError struct {
	ID any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Row with id %v not found", e.ID)
}
