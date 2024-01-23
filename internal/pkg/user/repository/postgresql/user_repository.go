package postgresql

import (
	"encoding/hex"
	"fmt"

	"github.com/yarikTri/web-transport-cards/internal/models"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

const UNSAFE_HARDCODE_SALT = "34xj&^R#*%&FTE^CWYUGhj"

type PostgreSQL struct {
	db *gorm.DB
}

func NewPostgreSQL(db *gorm.DB) *PostgreSQL {
	return &PostgreSQL{
		db: db,
	}
}

func (p *PostgreSQL) GetByID(userID int) (models.User, error) {
	var user models.User
	if err := p.db.First(&user, userID).Error; err != nil {
		return models.User{}, nil
	}

	return user, nil
}

func (p *PostgreSQL) GetByUsername(username string) (models.User, error) {
	var user models.User
	if err := p.db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, nil
	}

	return user, nil
}

func (p *PostgreSQL) GetByCreds(username, password string) (models.User, error) {
	var user models.User
	if err := p.db.Where("username = ? AND password = ?", username, p.GetSaltedHash(password)).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (p *PostgreSQL) Create(user models.User) (models.User, error) {
	user.Password = p.GetSaltedHash(user.Password)

	fmt.Println(user.Password)

	if err := p.db.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (p *PostgreSQL) CheckUserIsModerator(userID int) (bool, error) {
	user, err := p.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.IsModerator, nil
}

func (p *PostgreSQL) GetSaltedHash(password string) string {
	hashedBytePassword := argon2.IDKey([]byte(password), []byte(UNSAFE_HARDCODE_SALT), 3, 32*1024, 4, 32)
	return hex.EncodeToString(hashedBytePassword)
}
