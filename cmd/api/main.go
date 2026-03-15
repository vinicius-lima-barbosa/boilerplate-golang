package main

import (
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/server"
)

func main() {
	app := server.Start()

	app.Listen(":3000")
}
