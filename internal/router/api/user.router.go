package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/handler"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/middleware"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/repositories"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/services"
)

func NewUserRouter(app *fiber.App) *UserRouter {
	userRepo := repositories.NewUserRepository()
	usersService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	userHandler := handler.NewUserHandler(usersService)
	jwtMiddleware := middleware.JWTAuthMiddleware(authService)

	return &UserRouter{
		app:         app,
		userHandler: userHandler,
		jwtAuth:     jwtMiddleware,
	}
}

type UserRouter struct {
	app         *fiber.App
	userHandler handler.UserHandler
	jwtAuth     fiber.Handler
}

func (r *UserRouter) Setup(api fiber.Router) {
	users := api.Group("/users")
	users.Get("/get", r.jwtAuth, r.userHandler.GetUsers)
	users.Get("/get/:id", r.jwtAuth, r.userHandler.GetUserByID)
	users.Post("/create", r.jwtAuth, middleware.DBTransactionMiddleware(), r.userHandler.CreateUser)
	users.Put("/update/:id", r.jwtAuth, middleware.DBTransactionMiddleware(), r.userHandler.UpdateUser)
	users.Delete("/delete/:id", r.jwtAuth, middleware.DBTransactionMiddleware(), r.userHandler.DeleteUser)
}
