package server

import (
	config "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/database"

	"github.com/gofiber/swagger"
	_ "github.com/gofiber/swagger"

	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/router"

	"github.com/gofiber/fiber/v2"
)

func Start() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Boilerplate Golang v1.0.1",
	})

	cfg := swagger.Config{
		URL:       "/",
		ConfigURL: "/docs",
		Title:     "Documentação API Boilerplate Golang",
	}

	app.Use(swagger.New(cfg))
	config.Connection()

	router.Setup(app)

	return app
}
