package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yarikTri/web-transport-cards/internal/models"
)

type Redis struct {
	db  *redis.Client
	ctx context.Context
}

func NewRedis(rc *redis.Client) *Redis {
	return &Redis{
		db:  rc,
		ctx: context.Background(),
	}
}

func (r *Redis) GetTicketDraft(userID int) (models.Ticket, error) {
	var ticket models.Ticket

	err := r.db.Get(r.ctx, fmt.Sprint(userID)).Scan(&ticket)
	if err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (r *Redis) SetTicketDraft(ticket models.Ticket) (models.Ticket, error) {
	err := r.db.Set(r.ctx, fmt.Sprint(ticket.CreatorID), ticket, time.Duration((1>>32)*time.Hour)).Err()
	if err != nil {
		return models.Ticket{}, err
	}

	return r.GetTicketDraft(ticket.CreatorID)
}

func (r *Redis) DelTicketDraft(userID int) error {
	return r.db.Del(r.ctx, fmt.Sprint(userID)).Err()
}

func (r *Redis) AddRoute(userID int, route models.Route) (models.Ticket, error) {
	ticket, err := r.GetTicketDraft(userID)
	if err != nil {
		return models.Ticket{}, err
	}

	ticket.Routes = append(ticket.Routes, route)
	return r.SetTicketDraft(ticket)
}

func (r *Redis) DeleteRoute(userID int, routeID int) (models.Ticket, error) {
	ticket, err := r.GetTicketDraft(userID)
	if err != nil {
		return models.Ticket{}, err
	}

	currRoutes := ticket.Routes
	for ind, _route := range currRoutes {
		if int(_route.ID) == routeID {
			ticket.Routes = append(currRoutes[:ind], currRoutes[ind+1:]...)
			return r.SetTicketDraft(ticket)
		}
	}

	return ticket, nil
}
