package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	gorm.Model

	Name            string
	Active          bool `gorm:"default:true"`
	Capacity        uint32
	StartStation    string
	EndStation      string
	StartTime       string
	EndTime         string
	IntervalMinutes uint32
	Description     string
	ImageUUID       uuid.UUID
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
		ImageUUID:       r.ImageUUID,
	}
}

type RouteTransfer struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Active          bool      `json:"active"`
	Capacity        uint32    `json:"capacity"`
	StartStation    string    `json:"start_station"`
	EndStation      string    `json:"end_station"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	IntervalMinutes uint32    `json:"interval_minutes"`
	Description     string    `json:"description"`
	ImageUUID       uuid.UUID `json:"image_uuid"`
}
