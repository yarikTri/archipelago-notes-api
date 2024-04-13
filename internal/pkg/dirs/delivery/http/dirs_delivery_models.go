package http

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
)

type CreateDirRequest struct {
	Name        string `json:"name" valid:"required"`
	ParentDirID int    `json:"parent_dir_id"`
}

func (cdr *CreateDirRequest) validate() error {
	_, err := valid.ValidateStruct(cdr)
	return err
}

type UpdateDirRequest struct {
	ID   int    `json:"id" valid:"required"`
	Name string `json:"name" valid:"required"`
	Path string `json:"subpath" valid:"required"`
}

func (udr *UpdateDirRequest) validate() error {
	_, err := valid.ValidateStruct(udr)
	return err
}

func (udr *UpdateDirRequest) ToDir() models.Dir {
	return models.Dir{
		ID:   udr.ID,
		Name: udr.Name,
		Path: udr.Path,
	}
}
