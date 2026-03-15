package router

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app *fiber.App
}

func New(app *fiber.App) *Router {
	return &Router{
		app: app,
	}
}

func Setup(app *fiber.App) {
	router := New(app)
	app.Stack()
	// Setup API routes with rate limiter
	// api := app.Group("/api", limiter.New())

	// Setup API routes
	router.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
