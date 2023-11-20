package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Routes     []Route `gorm:"many2many:ticket_routes;"`
	State      string
	FormDate   string
	FinishDate string
	UserID     int
	User       User
}
