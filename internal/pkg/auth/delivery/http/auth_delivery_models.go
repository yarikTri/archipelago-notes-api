package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/yarikTri/web-transport-cards/internal/models"
)

type SignUpRequest struct {
	Username string `json:"username" valid:"required"`
	FullName string `json:"full_name" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (sur *SignUpRequest) validate() error {
	_, err := valid.ValidateStruct(sur)
	return err
}

func (sur *SignUpRequest) toUser() models.User {
	return models.User{
		Username:    sur.Username,
		FullName:    sur.FullName,
		Password:    sur.Password,
		IsModerator: false,
	}
}

type LoginRequest struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (sur *LoginRequest) validate() error {
	_, err := valid.ValidateStruct(sur)
	return err
}
