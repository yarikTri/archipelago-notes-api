package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"gorm.io/gorm"
)

type CreateRouteRequest struct {
	Name            string `json:"name" valid:"required"`
	Capacity        uint32 `json:"capacity" valid:"required"`
	StartStation    string `json:"start_station" valid:"required"`
	EndStation      string `json:"end_station" valid:"required"`
	StartTime       string `json:"start_time" valid:"required"`
	EndTime         string `json:"end_time" valid:"required"`
	IntervalMinutes uint32 `json:"interval_minutes" valid:"required"`
	Description     string `json:"description" valid:"required"`
}

func (crr *CreateRouteRequest) validate() error {
	_, err := valid.ValidateStruct(crr)
	return err
}

func (crr *CreateRouteRequest) ToRoute() models.Route {
	return models.Route{
		Name:            crr.Name,
		Capacity:        crr.Capacity,
		StartStation:    crr.StartStation,
		EndStation:      crr.EndStation,
		StartTime:       crr.StartTime,
		EndTime:         crr.EndTime,
		IntervalMinutes: crr.IntervalMinutes,
		Description:     crr.Description,
	}
}

type UpdateRouteRequest struct {
	Name            string `json:"name" valid:"required"`
	Capacity        uint32 `json:"capacity" valid:"required"`
	StartStation    string `json:"start_station" valid:"required"`
	EndStation      string `json:"end_station" valid:"required"`
	StartTime       string `json:"start_time" valid:"required"`
	EndTime         string `json:"end_time" valid:"required"`
	IntervalMinutes uint32 `json:"interval_minutes" valid:"required"`
	Description     string `json:"description" valid:"required"`
}

func (urr *UpdateRouteRequest) validate() error {
	_, err := valid.ValidateStruct(urr)
	return err
}

func (urr *UpdateRouteRequest) ToRoute(id uint64) models.Route {
	return models.Route{
		Model:           gorm.Model{ID: uint(id)},
		Name:            urr.Name,
		Capacity:        urr.Capacity,
		StartStation:    urr.StartStation,
		EndStation:      urr.EndStation,
		StartTime:       urr.StartTime,
		EndTime:         urr.EndTime,
		IntervalMinutes: urr.IntervalMinutes,
		Description:     urr.Description,
	}
}

type ListRoutesResponse struct {
	TicketDraftID *int                   `json:"ticket_draft_id"`
	Routes        []models.RouteTransfer `json:"routes"`
}
