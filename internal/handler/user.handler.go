package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/middleware"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/requests"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/responses"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/services"
	"gorm.io/gorm"
)

type UserHandler interface {
	GetUsers(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	// GetUserByID(c *fiber.Ctx) error
	// UpdateUser(c *fiber.Ctx) error
	// DeleteUser(c *fiber.Ctx) error
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags Users
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users [get]
func (h *userHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetUsers()
	if err != nil {
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return responses.NewSuccessResponseWithData(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User to create"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users [post]
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	tx, ok := c.Locals(middleware.DBTransactionKey).(*gorm.DB)
	if !ok || tx == nil {
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Database transaction not initialized", nil)
	}

	data := new(requests.CreateUserRequest)
	if err := c.BodyParser(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	userModel := data.ToModel()
	user, err := h.userService.WithTrx(tx).CreateUser(userModel)
	if err != nil {
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user", err)
	}

	return responses.NewSuccessResponseWithData(c, fiber.StatusCreated, "User created successfully", user)
}
