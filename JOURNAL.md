
TODO:

- [x] inside `/src/internal/`, write preliminary build files for _SQLc_
      __without__ running _SQLc_ for now.
- [ ] Read Chapter 3 in the HTMX Bible to find out about requirements
      for pages and templates.
- [ ] build SQLc system; may be tricky
- [ ] copy `db.Setup()` from Obsidian into project and into templates
- [ ] inside `/src/views/`, write the templates for the implementation 
      of the chapter 3 project.
- [ ] implement the project as a fiber project. May be tricky because of paths
      and templates.

## 2024-02-11 17:57

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

This has consequences. When SQLc returns aA `Null*` type property, from a database request, it does it like this:

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
    }                   b  
]
```

This is a simple JSON dump from a `SELECT` query. If a value is not `NULL` in
the database, in the model its `Valid` property will be `true`, otherwise
`false`.

For us this means that we have to write a helper function that converts this
result to a result we want to work with down the road.


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
