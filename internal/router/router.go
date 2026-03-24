package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/router/api"
)

type Router struct {
	app    *fiber.App
	health *api.HealthRouter
	users  *api.UserRouter
}

func New(app *fiber.App) *Router {
	return &Router{
		app:    app,
		health: api.NewHealthRouter(app),
		users:  api.NewUserRouter(app),
	}
}

func Setup(app *fiber.App) {
	router := New(app)
	app.Stack()
	// Setup API routes with rate limiter
	api := app.Group("/api", limiter.New())

	// Setup API routes
	router.health.Setup(api)
	router.users.Setup(api)
}
