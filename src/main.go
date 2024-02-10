package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", GetContacts)
	app.Get("/contacts", GetContacts)

	app.Listen(":5000")
}

func GetContacts(c *fiber.Ctx) error {
	searchTerm := c.Query("q")

	if searchTerm == "" {
		// => return all contacts as a list
		return c.SendString("No search string was given.")
	}

	// => return only the found contacts as a list
	return c.SendString(
		fmt.Sprintf("We have a search string: '%s'", searchTerm))
}
