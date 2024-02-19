package main

import (
	"context"
	"fmt"

	"database/sql"

	"github.com/Mt-Lampert/htmx_on_fiber/src/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/usepzaka/fiberflash"
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
	app.Post("/contacts/new", AddContact)

	app.Get("/contacts/:id", SingleContact)

	app.Get("/contacts/:id/edit", EditContact)

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
				"message": "Nothing Found",
			})
		}
		contacts := getProperContacts(rawContacts)
		return c.Render("pages/contacts", fiber.Map{
			"contacts": contacts,
			"query":    searchTerm,
			"Flash":    fiberflash.Get(c),
		}, "layouts/_baseof")
	}

	// => return only the found contacts as a list
	return c.SendString(
		fmt.Sprintf("We have a search string: '%s'", searchTerm))
}

func NewContact(c *fiber.Ctx) error {
	return c.Render("pages/contact-form", fiber.Map{
		"Data": fiber.Map{
			"Email": "charlie.cotton@cotton-charlie.com",
			"First": "Charlie",
			"Last":  "Cotton",
			"Phone": "1-58587193-8199",
		},
		"Flash": fiberflash.Get(c),
	}, "layouts/_baseof")
}

func AddContact(c *fiber.Ctx) error {
	ctx := context.Background()
	mp := fiber.Map{
		"Status": "success",
		"Msg":    "The new contact has been saved!",
	}
	// get the form data
	dbParams := db.AddContactParams{
		FirstName: sql.NullString{
			Valid:  true,
			String: c.FormValue("first"),
		},
		LastName: sql.NullString{
			Valid:  true,
			String: c.FormValue("last"),
		},
		Email: sql.NullString{
			Valid:  true,
			String: c.FormValue("email"),
		},
		Phone: sql.NullString{
			Valid:  true,
			String: c.FormValue("phone"),
		},
	}

	_, err := db.Qs.AddContact(ctx, dbParams)
	if err != nil {
		mp = fiber.Map{
			"Status": "error",
			"Msg":    "Could not save the new contact!",
		}
		return fiberflash.WithError(c, mp).Redirect("/contacts/new")
	}

	return fiberflash.WithSuccess(c, mp).Redirect("/contacts")
}

func SingleContact(c *fiber.Ctx) error {
	ctx := context.Background()
	id, _ := c.ParamsInt("id")
	scArg := sql.NullInt64{
		Valid: true,
		Int64: int64(id),
	}

	rawContact, err := db.Qs.GetContact(ctx, scArg)
	if err != nil {
		mp := fiber.Map{
			"Status": "error",
			"Msg":    fmt.Sprintf("Eh!!? Could not find a contact with id '%d'", id),
		}
		fiberflash.WithError(c, mp).Redirect("/contacts")
	}

	return c.Render(
		"pages/single-contact",
		getProperContact(rawContact),
		"layouts/_baseof",
	)
}

func EditContact(c *fiber.Ctx) error {
	ctx := context.Background()
	id, _ := c.ParamsInt("id")
	scArg := sql.NullInt64{
		Valid: true,
		Int64: int64(id),
	}

	rawContact, err := db.Qs.GetContact(ctx, scArg)
	if err != nil {
		mp := fiber.Map{
			"Status": "error",
			"Msg":    fmt.Sprintf("Eh!!? Could not find a contact with id '%d' for editing", id),
		}
		fiberflash.WithError(c, mp).Redirect("/contacts")
	}

	return c.Render("pages/contact-form", fiber.Map{
		"Data":  getProperContact(rawContact),
		"Flash": fiber.Map{},
	}, "layouts/_baseof")
}

// vim: foldmethod=indent
