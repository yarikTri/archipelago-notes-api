package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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

func (r *Redis) CreateSession(sessionID string, userID int, duration time.Duration) error {
	return r.db.Set(r.ctx, sessionID, fmt.Sprint(userID), duration).Err()
}

func (r *Redis) DeleteSession(sessionID string) error {
	return r.db.Del(r.ctx, sessionID).Err()
}

func (r *Redis) GetValueBySessionID(sessionID string) (int, error) {
	var userID int
	if err := r.db.Get(r.ctx, sessionID).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}
