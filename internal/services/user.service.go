package services

import (
	"context"

	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/dto"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/models"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/repositories"
	"gorm.io/gorm"
)

type UserService interface {
	GetUsers() ([]dto.UserResponse, error)
	CreateUser(user *models.User) (*dto.UserResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	UpdateUser(id uint, user *models.User) (*dto.UserResponse, error)
	DeleteUser(id uint) error
	WithTrx(trx *gorm.DB) UserService
}

func (s *userService) WithTrx(trx *gorm.DB) UserService {
	if trx == nil {
		return s
	}

	return &userService{
		userRepo: s.userRepo.WithTrx(trx),
	}
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	return dto.ToUserListResponse(users), nil
}

func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) CreateUser(user *models.User) (*dto.UserResponse, error) {
	createdUser, err := s.userRepo.Create(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(createdUser), nil
}

func (s *userService) UpdateUser(id uint, user *models.User) (*dto.UserResponse, error) {
	updatedUser, err := s.userRepo.Update(context.Background(), id, user)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(updatedUser), nil
}

func (s *userService) DeleteUser(id uint) error {
	user, err := s.userRepo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}

	return s.userRepo.Delete(context.Background(), id)
}
