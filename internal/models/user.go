package models

import "github.com/gofrs/uuid/v5"

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	RootDirID int       `db:"root_dir_id"`
}

func (u *User) ToTransfer() *UserTransfer {
	return &UserTransfer{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		RootDirID: u.RootDirID,
	}
}

type UserTransfer struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	RootDirID int    `json:"root_dir_id"`
}
