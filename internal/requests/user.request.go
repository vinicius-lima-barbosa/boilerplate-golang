package requests

import "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"

type CreateUserRequest struct {
	Name        string      `json:"name" validate:"required"`
	Email       string      `json:"email" validate:"required,email"`
	Type        models.Type `json:"type" validate:"required,oneof=admin user guest"`
	Description string      `json:"description"`
}

func (c *CreateUserRequest) ToModel() *models.User {
	return &models.User{
		Name:        c.Name,
		Email:       c.Email,
		Type:        c.Type,
		Description: c.Description,
	}
}
