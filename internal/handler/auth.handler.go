package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/middleware"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/requests"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/responses"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/services"
	"gorm.io/gorm"
)

type AuthHandler interface {
	Signup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type authHandler struct {
	authService services.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *authHandler) Signup(c *fiber.Ctx) error {
	tx, ok := c.Locals(middleware.DBTransactionKey).(*gorm.DB)
	if !ok || tx == nil {
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Database transaction not initialized", nil)
	}

	data := new(requests.SignupRequest)
	if err := c.BodyParser(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.validator.Struct(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid signup data", err)
	}

	result, err := h.authService.WithTrx(tx).Signup(data.Name, data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyInUse) {
			return responses.NewErrorResponse(c, fiber.StatusConflict, "Email already in use", nil)
		}
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Failed to signup", err)
	}

	return responses.NewSuccessResponseWithData(c, fiber.StatusCreated, "User signed up successfully", result)
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	tx, ok := c.Locals(middleware.DBTransactionKey).(*gorm.DB)
	if !ok || tx == nil {
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Database transaction not initialized", nil)
	}

	data := new(requests.LoginRequest)
	if err := c.BodyParser(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.validator.Struct(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid login data", err)
	}

	result, err := h.authService.WithTrx(tx).Login(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return responses.NewErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
		}
		if errors.Is(err, services.ErrInactiveUser) {
			return responses.NewErrorResponse(c, fiber.StatusForbidden, "Inactive user", nil)
		}
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Failed to login", err)
	}

	return responses.NewSuccessResponseWithData(c, fiber.StatusOK, "Login successful", result)
}

func (h *authHandler) Refresh(c *fiber.Ctx) error {
	tx, ok := c.Locals(middleware.DBTransactionKey).(*gorm.DB)
	if !ok || tx == nil {
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Database transaction not initialized", nil)
	}

	data := new(requests.RefreshRequest)
	if err := c.BodyParser(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.validator.Struct(data); err != nil {
		return responses.NewErrorResponse(c, fiber.StatusBadRequest, "Invalid refresh data", err)
	}

	result, err := h.authService.WithTrx(tx).Refresh(data.RefreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			return responses.NewErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired refresh token", nil)
		}
		if errors.Is(err, services.ErrInactiveUser) {
			return responses.NewErrorResponse(c, fiber.StatusForbidden, "Inactive user", nil)
		}
		return responses.NewErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", err)
	}

	return responses.NewSuccessResponseWithData(c, fiber.StatusOK, "Token refreshed successfully", result)
}
