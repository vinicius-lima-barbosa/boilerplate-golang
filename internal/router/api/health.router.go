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

// healthCheck godoc
// @Summary Verifica saude da API
// @Description Endpoint para validar se a API esta em execucao
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (r *HealthRouter) healthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status: "ok",
	})
}
