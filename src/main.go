package main

import (
	"github.com/Mt-Lampert/htmx_on_fiber/src/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".go.html")
	db.Setup()

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./static")

	build_routing(app)

	app.Listen(":5000")
}

// vim: foldmethod=indent
