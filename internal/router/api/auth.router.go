package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/handler"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/middleware"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/repositories"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/services"
)

type AuthRouter struct {
	app         *fiber.App
	authHandler handler.AuthHandler
}

func NewAuthRouter(app *fiber.App) *AuthRouter {
	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	return &AuthRouter{
		app:         app,
		authHandler: authHandler,
	}
}

func (r *AuthRouter) Setup(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Post("/signup", middleware.DBTransactionMiddleware(), r.authHandler.Signup)
	auth.Post("/login", middleware.DBTransactionMiddleware(), r.authHandler.Login)
	auth.Post("/refresh", middleware.DBTransactionMiddleware(), r.authHandler.Refresh)
}
