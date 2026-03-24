package dto

import (
	"time"

	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"
)

type UserResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Type        models.Type `json:"type"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Type:        user.Type,
		Description: user.Description,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func ToUserListResponse(users []models.User) []UserResponse {
	responses := make([]UserResponse, len(users))

	for i, user := range users {
		responses[i] = *ToUserResponse(&user)
	}

	return responses
}
