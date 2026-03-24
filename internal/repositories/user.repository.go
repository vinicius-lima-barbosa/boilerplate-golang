package repositories

import (
	"context"

	config "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/database"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, params *models.User) (*models.User, error)
	Update(ctx context.Context, id uint, params *models.User) (*models.User, error)
	Delete(ctx context.Context, id uint) error
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

func (u *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := u.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (u *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) Update(ctx context.Context, id uint, user *models.User) (*models.User, error) {
	var existingUser models.User
	err := u.db.WithContext(ctx).First(&existingUser, id).Error
	if err != nil {
		return nil, err
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Type = user.Type
	existingUser.Description = user.Description

	err = u.db.WithContext(ctx).Save(&existingUser).Error
	if err != nil {
		return nil, err
	}

	return &existingUser, nil
}

func (u *userRepository) Delete(ctx context.Context, id uint) error {
	return u.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
