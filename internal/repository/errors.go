package repository

import "fmt"

type NotFoundError struct {
	ID interface{}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("entity with ID %v not found", e.ID)
}
