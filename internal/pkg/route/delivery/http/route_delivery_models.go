package http

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"gorm.io/gorm"
)

type CreateRouteRequest struct {
	Name            string `json:"name" valid:"required"`
	Capacity        uint32 `json:"capacity" valid:"required"`
	StartStation    string `json:"start_station" valid:"required"`
	EndStation      string `json:"end_station" valid:"required"`
	StartTime       int64  `json:"start_time" valid:"required"`
	EndTime         int64  `json:"end_time" valid:"required"`
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
		StartTime:       time.Unix(crr.StartTime, 0),
		EndTime:         time.Unix(crr.EndTime, 0),
		IntervalMinutes: crr.IntervalMinutes,
		Description:     crr.Description,
	}
}

type UpdateRouteRequest struct {
	Name            string `json:"name" valid:"required"`
	Capacity        uint32 `json:"capacity" valid:"required"`
	StartStation    string `json:"start_station" valid:"required"`
	EndStation      string `json:"end_station" valid:"required"`
	StartTime       int64  `json:"start_time" valid:"required"`
	EndTime         int64  `json:"end_time" valid:"required"`
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
		StartTime:       time.Unix(urr.StartTime, 0),
		EndTime:         time.Unix(urr.EndTime, 0),
		IntervalMinutes: urr.IntervalMinutes,
		Description:     urr.Description,
	}
}
