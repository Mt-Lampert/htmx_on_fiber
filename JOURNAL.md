
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
