package usecase

import (
	"time"

	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/auth"
	"github.com/yarikTri/web-transport-cards/internal/pkg/user"
)

// Usecase implements auth.Usecase
type Usecase struct {
	authRepo auth.Repository
	userRepo user.Repository
}

func NewUsecase(ar auth.Repository, ur user.Repository) *Usecase {
	return &Usecase{
		authRepo: ar,
		userRepo: ur,
	}
}

func (u *Usecase) GetUserBySessionID(sessionID string) (models.User, error) {
	userID, err := u.authRepo.GetValueBySessionID(sessionID)
	if err != nil {
		return models.User{}, err
	}

	return u.userRepo.GetByID(userID)
}

func (u *Usecase) CheckUserIsModerator(userID int) (bool, error) {
	return u.userRepo.CheckUserIsModerator(userID)
}

func (u *Usecase) SignUp(user models.User) (models.User, error) {
	return u.userRepo.Create(user)
}

func (u *Usecase) Login(username, password string, sessionDuration time.Duration) (string, error) {
	user, err := u.userRepo.GetByCreds(username, password)
	if err != nil {
		return "", err
	}

	sessionID := u.userRepo.GetSaltedHash(username)
	u.authRepo.CreateSession(sessionID, int(user.ID), sessionDuration)

	return sessionID, err
}

func (u *Usecase) Logout(sessionID string) error {
	return u.authRepo.DeleteSession(sessionID)
}
