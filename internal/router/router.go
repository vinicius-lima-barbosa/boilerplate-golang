package router

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app *fiber.App
}

type HealthResponse struct {
	Status string `json:"status"`
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
	router.app.Get("/api/health", healthCheck)
}

// healthCheck godoc
// @Summary Verifica saude da API
// @Description Endpoint para validar se a API esta em execucao
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/health [get]
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{Status: "ok"})
}
