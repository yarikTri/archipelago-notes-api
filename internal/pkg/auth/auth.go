package auth

import (
	"time"

	"github.com/yarikTri/web-transport-cards/internal/models"
)

type Usecase interface {
	GetUserBySessionID(sessionID string) (models.User, error)
	CheckUserIsModerator(userID int) (bool, error)
	SignUp(user models.User) (models.User, error)
	Login(username, password string, sessionDuration time.Duration) (sessionID string, err error)
	Logout(sessionID string) error
}

type Repository interface {
	CreateSession(sessionID string, userID int, duration time.Duration) error
	DeleteSession(sessionID string) error
	GetValueBySessionID(sessionID string) (int, error)
}
