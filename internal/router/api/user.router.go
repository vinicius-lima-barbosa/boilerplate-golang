package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/handler"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/middleware"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/repositories"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/services"
)

type UserRouter struct {
	app         *fiber.App
	userHandler handler.UserHandler
}

func NewUserRouter(app *fiber.App) *UserRouter {
	userRepo := repositories.NewUserRepository()
	usersService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(usersService)

	return &UserRouter{
		app:         app,
		userHandler: userHandler,
	}
}

func (r *UserRouter) Setup(api fiber.Router) {
	users := api.Group("/users")
	users.Get("/get", r.userHandler.GetUsers)
	users.Post("/create", middleware.DBTransactionMiddleware(), r.userHandler.CreateUser)
	// r.app.Post("/users", r.userHandler.CreateUser)
	// r.app.Get("/users/:id", r.userHandler.GetUserByID)
	// r.app.Put("/users/:id", r.userHandler.UpdateUser)
	// r.app.Delete("/users/:id", r.userHandler.DeleteUser)
}
