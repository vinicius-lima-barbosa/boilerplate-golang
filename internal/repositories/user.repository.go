package repositories

import (
	"context"

	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) Create(ctx context.Context, params *models.User) (*models.User, error) {
	user := &models.User{
		Name:        params.Name,
		Email:       params.Email,
		Type:        params.Type,
		Description: params.Description,
	}

	err := u.db.WithContext(ctx).Create(&user).Error
	return user, err
}
