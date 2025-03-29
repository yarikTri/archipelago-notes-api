package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
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

func (sr *SessionsRepository) ClearAllSessions() error {
	// Get all keys
	keys, err := sr.db.Keys(context.TODO(), "*").Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	// If no keys found, return nil
	if len(keys) == 0 {
		return nil
	}

	// Delete all keys
	if err := sr.db.Del(context.TODO(), keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	return nil
}
