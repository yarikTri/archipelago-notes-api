package models

import "time"

type Ticket struct {
	ID        uint32
	RouteID   uint32
	State     string
	StartedAt time.Time
}
