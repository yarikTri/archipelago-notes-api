package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const (
	DRAFT_STATE    = "draft"
	DELETED_STATE  = "deleted"
	FORMED_STATE   = "formed"
	APPROVED_STATE = "approved"
	REJECTED_STATE = "rejected"
	ENDED_STATE    = "ended"

	DEFAULT_CREATOR_ID = 1
)

type Ticket struct {
	gorm.Model

	Routes      []Route `gorm:"many2many:ticket_routes;"`
	State       string  `gorm:"default:draft"`
	WriteState  *string
	FormTime    time.Time
	ApproveTime time.Time
	EndTime     time.Time
	CreatorID   int
	Creator     User `gorm:"foreignKey:CreatorID"`
	ModeratorID *int
	Moderator   *User `gorm:"foreignKey:ModeratorID"`
}

func (t Ticket) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Ticket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *Ticket) ToTransfer() TicketTransfer {
	routesTransfers := make([]RouteTransfer, 0)
	for _, route := range t.Routes {
		routesTransfers = append(routesTransfers, route.ToTransfer())
	}

	endTimeUnix := t.EndTime.Unix()
	var endTime *int64 = &endTimeUnix
	if *endTime == -62135596800 {
		endTime = nil
	}

	formTimeUnix := t.FormTime.Unix()
	var formTime *int64 = &formTimeUnix
	if *formTime == -62135596800 {
		formTime = nil
	}

	return TicketTransfer{
		ID:              t.ID,
		Routes:          routesTransfers,
		State:           t.State,
		WriteState:      t.WriteState,
		CreateTime:      t.CreatedAt.Unix(),
		FormTime:        formTime,
		EndTime:         endTime,
		CreatorUsername: t.Creator.Username,
		ModeratorID:     t.ModeratorID,
	}
}

type TicketTransfer struct {
	ID              uint            `json:"id"`
	Routes          []RouteTransfer `json:"routes"`
	State           string          `json:"state"`
	WriteState      *string         `json:"write_state"`
	CreateTime      int64           `json:"create_time"`
	FormTime        *int64          `json:"form_time"`
	EndTime         *int64          `json:"end_time"`
	CreatorUsername string          `json:"creator_username"`
	ModeratorID     *int            `json:"moderator_id"`
}
