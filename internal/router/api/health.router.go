package api

import "github.com/gofiber/fiber/v2"

type HealthRouter struct {
	app *fiber.App
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHealthRouter(app *fiber.App) *HealthRouter {
	return &HealthRouter{
		app: app,
	}
}

func (r *HealthRouter) Setup(api fiber.Router) {
	api.Get("/health", r.healthCheck)
}

func (r *HealthRouter) healthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status: "ok",
	})
}
