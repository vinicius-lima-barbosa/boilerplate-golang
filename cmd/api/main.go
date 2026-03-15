package main

import (
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/server"
)

// @title Boilerplate Golang API
// @version v1.0.1
// @description API boilerplate em Golang utilizando Fiber, GORM e PostgreSQL
// @host localhost:3000
// @BasePath /

func main() {
	app := server.Start()

	app.Listen(":3000")
}
