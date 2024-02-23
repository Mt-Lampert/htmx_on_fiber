package main

import "github.com/gofiber/fiber/v2"

func build_routing(app *fiber.App) {
	app.Get("/", GetContacts)
	app.Get("/contacts", GetContacts)

	app.Get("/contacts/new", NewContact)
	app.Post("/contacts/new", AddContact)

	app.Get("/contacts/:id", SingleContact)

	app.Get("/contacts/:id/edit", EditContact)
	app.Post("contacts/:id/edit", UpdateContact)

	// app.Post("contacts/:id/delete", DeleteContact)
	app.Delete("/contacts/:id", DeleteContact)
}
