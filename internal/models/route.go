package models

import "gorm.io/gorm"

type Route struct {
	gorm.Model
	Name string
	// Tickets         []*Ticket `gorm:"many2many:route_ticket;"`
	Active          bool
	Capacity        uint32
	StartStation    string
	EndStation      string
	StartTime       string
	EndTime         string
	IntervalMinutes uint32
	Description     string
	ImagePath       *string
}

func (r *Route) ToTransfer() RouteTransfer {
	return RouteTransfer{
		ID:              r.ID,
		Name:            r.Name,
		Active:          r.Active,
		Capacity:        r.Capacity,
		StartStation:    r.StartStation,
		EndStation:      r.EndStation,
		IntervalMinutes: r.IntervalMinutes,
		StartTime:       r.StartTime,
		EndTime:         r.EndTime,
		Description:     r.Description,
		ImagePath:       r.ImagePath,
	}
}

type RouteTransfer struct {
	ID              uint
	Name            string
	Active          bool
	Capacity        uint32
	StartStation    string
	EndStation      string
	StartTime       string
	EndTime         string
	IntervalMinutes uint32
	Description     string
	ImagePath       *string
}
