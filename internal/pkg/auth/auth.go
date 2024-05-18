package auth

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type Usecase interface {
	GetUserIDBySessionID(sessionID string) (uuid.UUID, error)
	SignUp(email, name, password string) (uuid.UUID, error)
	Login(email, password string) (string, time.Duration, error)
	Logout(sessionID string) error
}

type SessionsRepository interface {
	GetUserIDBySessionID(sessionID string) (uuid.UUID, error)
	CreateSession(sessionID string, userID uuid.UUID, expiration time.Duration) error
	DeleteSession(sessionID string) error
}

type UsersRepository interface {
	GetUserIDAndPasswordByEmail(email string) (uuid.UUID, string, error)
	CreateUser(email, name, passwordHash string) (uuid.UUID, error)
}
