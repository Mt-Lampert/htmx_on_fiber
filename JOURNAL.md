
## TODO:



## 2024-02-23 05:36

- [x] Logger installieren
- [x] Logger mit neuem Endpoint ausprobieren; rumspielen, experimentieren.

Hab ihn installiert; er loggt jetzt alles. Und ich hab was herausgefunden:
`fmt.Printf()` lässt sich für alles spezielle Logging ausnutzen.

Die Arbeit mit HTMX ist ... „anders“. Mein größtes Problem ist es, in variablen
Frames zu denken. Beispiel: Ich mach eine Abfrage. Diese Abfrage hat ein
Ergebnis (⇒ snippet) und soll auch eine _Flash-Message_ absondern.

Mit HTMX müsste der Container für die Flash-Message in das Snippet eingebaut werden.
Im CSS wäre dann definiert, dass die Message dann `fixed` ist und z.B. oben
rechts aufleuchtet und nach ein paar Sekunden wieder verschwindet. Das kann man
ja alles so einstellen. Aber funktioniert es auch dann, wenn HTMX das Snippet
hineinlädt? Die Doku sagt Ja. wir müssen das einfach ausprobieren. Deshalb
folgender Plan:

1. CSS-Klassen für Flash-Nachrichten einbauen
0. `DELETE /contacts/:id` zu einem HX-Element umbauen.
0. `hx-swap='#contact-table'`
0. im Snippet den flash einbauen, mit CSS für Positionierung und Transitions
0. beten, dass alles richtig funktioniert.


## 2024-02-20 12:11

- [x] implement `DELETE /contacts/:id`

The easiest implementation of all, if it hadn't been for an error in
`src/internal/queries.sql`. This error prevented the scaffolding of the
`DeleteContact()` method; after correcting and running `$ sqlc generate` again,
everything turned out fine.

This also showed to me that my selfed-baked `init.go` file went totally
untouched by _SQLc_. “Most satisfactory!”, as Mr Wolfe used to say.

## 2024-02-20 10:42

- [x] implement `POST /contacts/:id/edit`

Quite straightforward. Had to consult the `db.Qs.UpdateContact()` code to do it
right, and finally it worked out.

However, reusing `views/pages/contact-form.go.html` had a nasty surprise in
store: Instead of updating the contact, it added another one using the data I
gave it. 

The problem was a hard-coded `/contacts/new` as the form action – the form
always ended up adding a new contact in the database. After adding the
`"Action"` prop to the `NewContact()/UpdateContact()` handlers, things are
working well.


## 2024-02-19 22:04

- [x] Read the final part of the HTMX Bible chapter
- [x] Add more well-formed TODOs
- [x] implement `GET /contacts/:id/edit`

Easy as pie. Just had to copy and adapt the database code from `SingleView()`
and send the result to `views/pages/contact-form.go.html`.

## 2024-02-19 18:12

- [x] implement `GET /contacts/:id`

Not too hard; could make use of the stuff I built before; but Go requires me to
be so goddamn explicit!

## 2024-02-19 18:12

- [x] implement `POST /contacts/new`, Flash part

The hardest thing was getting my head around the concept (see my Obisidian
article about it); after that, the implementation was a breeze.

__p.a.:__ Got a warning from Firefox that the current implementation of
fiberflash won't be tolerated much longer because of a cookie with  Same Site
issues. Let's see how this works out. At the moment, it just works :) 


## 2024-02-18 23:56

- [x] implement `POST /contacts/new`

To follow the model in the HTMX bible, I had to install the _fiberflash_
package, understand it and include it into the project.

I finally made it, and I faced a Go-typical surprise when I had to satisfy the
type requirements for `db.AddContactParams` and `database.sql.NullString`, just
to insert a few lousy data from the HTML form into the database. It's a pain in
the ..., yes indeed, but the code will always be workable by being so explicit.

## 2024-02-18 11:21

- [x] implement `GET /contacts/new`

Yes, I did it. After collecting some experience, this wasn't too tough. 

## 2024-02-11 19:47

- [x] inside `/src/views/`, write the templates for the implementation 
      of the chapter 3 project; steal the Markup from their project repository.
- [x] implement the project as a fiber project. May be tricky because of paths
      and templates.

I did it! Templates are working just as they should be. 

__Lessons Learned:__ In a template, defining a block with `{{ define }}` won't
have the block being rendered automatically. Only `{{ template }}` will do the
trick. It can even be added right below the definition!

```html
{{/* only defines it, doesn't show it! */}}
{{ define "my_block" }}
  <h1>Can you see me?</h1>
{{ end }}

{{/* THIS one shows it! */}}
{{ template "my_block" }}
```

## 2024-02-11 19:47

Another milestone. I built a conversion function to turn raw SQLc data into
something I can work with. It's a simple helper function I can now use whenever
I need proper database results.


## 2024-02-11 17:57

- [x] inside `/src/internal/`, write preliminary build files for _SQLc_
      __without__ running _SQLc_ for now.
- [x] build SQLc system; may be tricky
- [x] copy `db.Setup()` from Obsidian into project and into templates
- [x] test basic functionality of SQLc using a JSON dump.

I learned a lot again. SQLc works like a breeze, but it has its quirks. One of
them is the handling of `NULL` values in the database.

In case the _schema_ doesn't provide `NOT NULL` for a table property, SQLc will
assign an `sql.Null*` type in the model it creates for the table

```go
type Contact struct {
	ID        sql.NullInt64
	LastName  sql.NullString
	FirstName sql.NullString
	Phone     sql.NullString
	Email     sql.NullString
}
```

This has consequences. When SQLc returns aA `Null*` type property, from a
database request, it does it like this:

```json
[
    {
        "Email": {
            "String": "bob.bookie@bookiebob.com",
            "Valid": true
        },
        "FirstName": {
            "String": "Bob",
            "Valid": true
        },
        "ID": {
            "Int64": 2,
            "Valid": true
        },
        "LastName": {
            "String": "Booqie",
            "Valid": true
        },
        "Phone": {
            "String": "1-917-4890931",
            "Valid": true
        }
    }
]
```

This is a simple JSON dump from a `SELECT` query. If a value is not `NULL` in
the database, in the model its `Valid` property will be `true`, otherwise
`false`.

For us this means that we have to write a helper function that converts this
result to a result we want to work with down the road.

The [official Documentation](https://pkg.go.dev/database/sql#NullString) gives
a roadmap on what to do. Simple but tedious.

## 2024-02-11 11:20

I got SQLc up and running, that is, it initializes.

__Important Lesson:__ It matters a lot from where we invoke `go run`. So I will
add a task in the Makefile in the project root. If I run `make go_run`,
everything will be fine because all the paths will be set correctly here.

Running `:GoRun` from within Neovim with `main.go` open will move the `$PWD` to
`/src`, and Go complains it can't find files like `.env` any more.


## 2024-02-11 07:04

To go on with the Chapter, the following requirements must be met:

1. SQLc must be up and running, including the `Setup()` routine
0. SQLc must be tested, so there must be some data in the DB 
0. The SQLc model data type must be conforming to JSON and Form.


## 2024-02-08 19:01

These are the steps I took to get Go and Fiber up and running:

1. Install Go (Fedora 39):
    ```bash
    $ sudo dnf install golang-bin
    $ go version
    go version go1.21.6 linux/amd64
    ```
2. Set up the project as Golang project:
    ```bash
    $ go mod init github.com/mt-lampert/htmx_on_fiber
    ```

3. Include the _Fiber_ package:
   ```bash
   $ go get github.com/gofiber/fiber/v2
   ```

4. Test the following code
    ```go
    package main

    import "github.com/gofiber/fiber/v2"

    func main() {
        app := fiber.New()

        app.Get("/", func(c *fiber.Ctx) error {
            return c.SendString("Go Fiber is up and running!")
        })

        app.Listen(":3000")
    }
    ```

    I used [httpie](https://httpie.io/cli) to test the output:

    ```bash
    $ http :3000
    ```

Everything was fine!
