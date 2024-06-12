package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionIDLength = 32
	sessionTTL      = 90 * 24 * time.Hour
)

// Usecase implements auth.Usecase
type Usecase struct {
	sessionsRepo auth.SessionsRepository
	usersRepo    auth.UsersRepository
}

func NewUsecase(sr auth.SessionsRepository, ur auth.UsersRepository) *Usecase {
	return &Usecase{
		sessionsRepo: sr,
		usersRepo:    ur,
	}
}

func (u *Usecase) getPasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (u *Usecase) generateSessionID(length uint) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (u *Usecase) GetUserIDBySessionID(sessionID string) (uuid.UUID, error) {
	return u.sessionsRepo.GetUserIDBySessionID(sessionID)
}

func (u *Usecase) SignUp(email, name, password string) (string, uuid.UUID, time.Duration, error) {
	userID, err := u.usersRepo.CreateUser(email, name, u.getPasswordHash(password))
	if err != nil {
		return "", uuid.Max, 0, err
	}

	sessionID := u.generateSessionID(sessionIDLength)
	if err := u.sessionsRepo.CreateSession(sessionID, userID, sessionTTL); err != nil {

		if errDelete := u.usersRepo.DeleteUser(userID); errDelete != nil {
			// VERY SUS MOMENT
			err = fmt.Errorf("failed to delete user after create session err: %w", err)
		}

		return "", uuid.Max, 0, err
	}

	return sessionID, userID, sessionTTL, nil
}

func (u *Usecase) Login(email, password string) (string, uuid.UUID, time.Duration, error) {
	userID, passwordHash, err := u.usersRepo.GetUserIDAndPasswordByEmail(email)
	if err != nil {
		return "", uuid.Max, 0, err
	}

	gotPasswordHash := u.getPasswordHash(password)

	if passwordHash != gotPasswordHash {
		return "", uuid.Max, 0, fmt.Errorf("passwords dont match")
	}

	sessionID := u.generateSessionID(sessionIDLength)
	if err := u.sessionsRepo.CreateSession(sessionID, userID, sessionTTL); err != nil {
		return "", uuid.Max, 0, err
	}

	return sessionID, userID, sessionTTL, nil
}

func (u *Usecase) Logout(sessionID string) error {
	return u.sessionsRepo.DeleteSession(sessionID)
}
