package models

import "github.com/gofrs/uuid/v5"

type User struct {
	ID        uuid.UUID `db:"id"`
	Login     string    `db:"login"`
	Name      string    `db:"name"`
	RootDirID int       `db:"root_dir_id"`
}

func (u *User) ToTransfer() *UserTransfer {
	return &UserTransfer{
		ID:        u.ID.String(),
		Login:     u.Login,
		Name:      u.Name,
		RootDirID: u.RootDirID,
	}
}

type UserTransfer struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	RootDirID int    `json:"root_dir_id"`
}
