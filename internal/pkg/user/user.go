package user

import "github.com/yarikTri/web-transport-cards/internal/models"

type Repository interface {
	GetByUsername(username string) (models.User, error)
	GetByID(userID int) (models.User, error)
	CheckUserIsModerator(userID int) (bool, error)
	GetByCreds(username, password string) (models.User, error)
	Create(user models.User) (models.User, error)
	GetSaltedHash(password string) string
}
