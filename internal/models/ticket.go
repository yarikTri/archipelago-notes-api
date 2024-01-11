package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const (
	DRAFT_STATE     = "draft"
	DELETED_STATE   = "deleted"
	FORMED_STATE    = "formed"
	APPROVED_STATE  = "approved"
	REJECTED_STATE  = "rejected"
	ENDED_STATE     = "ended"
	FINALIZED_STATE = "finalized"

	DEFAULT_CREATOR_ID = 1
)

type Ticket struct {
	gorm.Model

	Routes      []Route `gorm:"many2many:ticket_routes;"`
	State       string  `gorm:"default:draft"`
	FormTime    time.Time
	ApproveTime time.Time
	EndTime     time.Time
	CreatorID   int  `gorm:"default:1"`
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

	if t.EndTime.Before(time.Now()) && int(t.EndTime.Unix()) != -62135596800 {
		t.State = ENDED_STATE
	}

	endTimeUnix := int(t.EndTime.Unix())
	var endTime *int = &endTimeUnix
	if int(t.EndTime.Unix()) == -62135596800 {
		endTime = nil
	}

	formTimeUnix := int(t.FormTime.Unix())
	var formTime *int = &formTimeUnix
	if int(t.EndTime.Unix()) == -62135596800 {
		formTime = nil
	}

	return TicketTransfer{
		ID:          t.ID,
		Routes:      routesTransfers,
		State:       t.State,
		CreateTime:  int(t.CreatedAt.Unix()),
		FormTime:    formTime,
		EndTime:     endTime,
		CreatorID:   t.CreatorID,
		ModeratorID: t.ModeratorID,
	}
}

type TicketTransfer struct {
	ID          uint            `json:"id"`
	Routes      []RouteTransfer `json:"routes"`
	State       string          `json:"state"`
	CreateTime  int             `json:"create_time"`
	FormTime    *int            `json:"form_time"`
	EndTime     *int            `json:"end_time"`
	CreatorID   int             `json:"creator_id"`
	ModeratorID *int            `json:"moderator_id"`
}
