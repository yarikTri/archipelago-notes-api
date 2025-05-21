package auth

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Usecase interface {
	GetUserIDBySessionID(sessionID string) (uuid.UUID, error)
	SignUp(email, name, password string) (string, uuid.UUID, time.Duration, error)
	Login(email, password string) (string, uuid.UUID, time.Duration, error)
	Logout(sessionID string) error
	ClearAllSessions() error
}

type SessionsRepository interface {
	GetUserIDBySessionID(sessionID string) (uuid.UUID, error)
	CreateSession(sessionID string, userID uuid.UUID, expiration time.Duration) error
	DeleteSession(sessionID string) error
	ClearAllSessions() error
}

type UsersRepository interface {
	GetUserIDAndPasswordByEmail(email string) (uuid.UUID, string, error)
	CreateUser(email, name, passwordHash string) (uuid.UUID, error)
	DeleteUser(userID uuid.UUID) error
}
