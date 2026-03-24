package repositories

import (
	"context"

	config "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/database"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, params *models.User) (*models.User, error)
	WithTrx(trx *gorm.DB) UserRepository
}

func (r *userRepository) WithTrx(trx *gorm.DB) UserRepository {
	if trx == nil {
		return r
	}

	return &userRepository{db: trx}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: config.DB.DB,
	}
}

func (u *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := u.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (u *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
