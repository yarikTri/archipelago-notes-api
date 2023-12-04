package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	FullName   string
	Username   string
	AvatarPath *string
}
