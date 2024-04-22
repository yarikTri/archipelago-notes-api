package http

import "github.com/yarikTri/archipelago-notes-api/internal/models"

type SearchUsersResponse struct {
	Users []*models.UserTransfer `json:"users"`
}
