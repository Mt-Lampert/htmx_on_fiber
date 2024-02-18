package main

import (
	"context"
	"fmt"

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

	app.Get("/", GetContacts)
	app.Get("/contacts", GetContacts)

	app.Get("/contacts/new", NewContact)

	app.Listen(":5000")
}

func GetContacts(c *fiber.Ctx) error {
	ctx := context.Background()
	searchTerm := c.Query("q")

	if searchTerm == "" {
		// => return all contacts as a list
		rawContacts, err := db.Qs.GetAllContacts(ctx)
		if err != nil || len(rawContacts) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "Not Found",
				"message": "Nothing Found"})
		}
		contacts := getProperContacts(rawContacts)
		return c.Render("pages/contacts", fiber.Map{
			"contacts": contacts,
			"query":    searchTerm,
		}, "layouts/_baseof")
	}

	// => return only the found contacts as a list
	return c.SendString(
		fmt.Sprintf("We have a search string: '%s'", searchTerm))
}

func NewContact(c *fiber.Ctx) error {

	return c.Render("pages/contact-form", fiber.Map{
		"Email": "charlie.cotton@cotton-charlie.com",
		"First": "Charlie",
		"Last":  "Cotton",
		"Phone": "1-58587193-8199",
	}, "layouts/_baseof")
}

// vim: foldmethod=indent
