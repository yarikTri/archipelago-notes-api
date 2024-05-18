package redis

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
	"time"
)

// SessionsRepository implements auth.SessionsRepository
type SessionsRepository struct {
	db *redis.Client
}

func NewSessionsRepository(db *redis.Client) *SessionsRepository {
	return &SessionsRepository{
		db: db,
	}
}

func (sr *SessionsRepository) GetUserIDBySessionID(sessionID string) (uuid.UUID, error) {
	return uuid.FromString(sr.db.Get(context.TODO(), sessionID).String())
}

func (sr *SessionsRepository) CreateSession(sessionID string, userID uuid.UUID, expiration time.Duration) error {
	return sr.db.Set(context.TODO(), sessionID, userID.String(), expiration).Err()
}

func (sr *SessionsRepository) DeleteSession(sessionID string) error {
	return sr.db.Del(context.TODO(), sessionID).Err()
}
