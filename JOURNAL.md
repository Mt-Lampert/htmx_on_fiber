
## TODO:

Nothing left at the moment

## 2024-02-23 19:21

- [x] implement auto-search feature
    
Auto-search means that you “automagically” get less results the more you type
into the search field. In this case, I want the auto-search mechanism be
triggered after I have typed at least four characters into the input field.

Yes! It finally worked! Here's the working code:

```html
<!-- file: views/pages/contacts.go.html -->
<div>
  <label for="search">search term</label>
  <input id="search" type="search" name="q" value="" placeholder="Search..."
         hx-get="/contacts/search"
         hx-trigger="keyup[this.value.length > 3] delay:200ms changed"
         hx-target="#theContacts"
  >
</div>
```

#### Annotations

1. The URL sent to the backend would be `/contacts/search?q={this.value}`.
0. `hx-trigger` Translates to 
    - trigger the GET request 
    - when there is a change
    - but only at a keyup event
    - and only when `value` has at least 4 characters
    - and then only after a delay of 200 milliseconds after the last keyup event
0. `hx-target` declares where the resulting HTML from the server has to be included.


### The handler

```go
func SearchContacts(c *fiber.Ctx) error {
	ctx := context.Background()
	searchTerm := fmt.Sprintf("%%%s%%", c.Query("q")) // -- 1 --

	theArgs := db.SearchContactsParams{
		FirstName: sql.NullString{Valid: true, String: searchTerm},
        // ...
	}

	rawContacts, err := db.Qs.SearchContacts(ctx, theArgs)
	if err != nil || len(rawContacts) == 0 {
		return c.SendString(fmt.Sprintf(
			"<p>Too bad! No contacts found matching '%s'</p>",
			searchTerm,
		))
	}
    // transform rawContacts into contacts we can use in the template
	contacts := getProperContacts(rawContacts)
	return c.Render("snippets/the_contacts", fiber.Map{
		"contacts": contacts,
	})
}
```

The code is quite self-explanatory, with one exception.

- `fmt.Sprintf("%%%s%%", c.Query("q"))`. If the `q` parameter was `ding`, then
  `fmt.Sprintf()` would return `"%ding%"` – which is exactly what is needed in
  an SQL query like this: `SELECT * FROM CONTACTS WHERE first_name LIKE
  '%ding%';`



__Lessons Learned:__ Stick to the details and read the error messages in the
browser! At the first tries I wrote `hx-target="theTable"`, not
`hx-target="#theTable"`. Since I did ___not___ read the error messages in the
browser, I lost half a day wondering what might be wrong with the HTMX version
I was using ...




## 2024-02-23 19:21

- [x] implement the Dropdown paging project (see
  [here](../ideas/pages_by_session.md) for details
  - [x] add more contacts to the database until it carries 22 contacts
  - [x] update `GetAllContacts` in `src/internal/queries.sql` according to the
    feature plans
  - [x] set a global variables `gSETS` and `gSET_SIZE` in `src/handlers.go`
  - [x] refactor the existing code in the `GetContacts()` to make sure it initially shows the first 10 contacts
  - [x] Add a 'More' button to add 5 contacts to the list and update the list accordingly
  - [x] Add a 'Reset' button to reset the display to 10 contacts

It went like a charm! The only thing I missed the first time was creating a
template snippet for the `More`/`Reset` button.



## 2024-02-23 19:21

- [x] Validation for unique mail.
    - [x] supplement `db` with specific mail search for user
    - [x] add HTMX to execute the request supplemented above

Yes, we made it! and we made it better than the original, and `sql.NullInt` made it even easier!

First, here's what we changed in the form:


```html
<p>
    <label for="email">Email</label> 
    <input name="email" id="email" type="email" placeholder="Email" value="{{ $email  }}"
      hx-get="/contacts/mailcheck"
      hx-target="next .error"
      hx-swap="outerHTML"
    >
    <span class="error">{{ .Errors.email }}</span>
</p>
```

#### Annotations:

It's essential to notice here that the HTMX is happening inside an `<input>`
element. This has consequences for the `hx-get` property:

1. As soon as the focus leaves this `<input>` element,  `hx-get` will send a
   'GET' request to `/contacts/mailcheck?{{ this.name }}={{ this.value }}` like 
   `/contacts/mailcheck?email=kim.kinky@weirdoclub.org`, for example.
   Why? Because of the default settings inside HTMX: If an `<input>` element
   sends `hx-get`, this will be what happens by default! Cool, eh?
2. `hx-target` will be the `next` DOM element with an 'error' class.
3. `hx-swap` will totally replace `hx-target` with the HTML response it
   receives from the backend.

That's it from the HTMX side. Now let's have a look at the handler:

```go
// file: /src/handlers.go
func CheckEmail(c *fiber.Ctx) error {
	ctx := context.Background()
	theMail := c.Query("email")  // -- 1 -- 
	rawID, err := db.Qs.GetEmail(ctx, sql.NullString{
		Valid:  true,
		String: theMail,
	}) // --2--

	if err != nil { // --3-- 
		return c.SendString("<span></span>")
	}

	if rawID.Valid {  // --4--
		return c.SendString("<span class='error'>Please choose another email!</span>")
	}

	return c.SendString("<span></span>") // --5--
}

// file: /src/internal/db/queries.sql.go
func (q *Queries) GetEmail(ctx context.Context, email sql.NullString) (sql.NullInt64, error) { ... }
```

First, I had to add a new entry in `/src/internal/queries.sql` to generate
`db.GetEmail()`. _SQLc_ then came up with the function shown here:

#### Annotations

The code gets the email from the query params (1), makes the database request
(2) by creating a valid `sql.NullString` argument and receives a result. 

Now the logic is important: If the database query was successful, this will be
taken for an __error!__ (4) Why? Because mail addresses must be unique in the
database. Everything else will be taken to be OK. (3,5)

__IMPORTANT:__ HTMX depends on the response to be valid HTML. That's why we
send an empty `<span>` element back and not an empty string. If this rule is
violated, HTMX will react drastically: it will take the entire landing page for
the response and, in this case, replace the `next .error` element with it
(`hx-target`)! 

So, if we suddenly find the landing page included in our current page, that's
the way HTMX tells us we made an HTML error ...



## 2024-02-23 15:51

- [x] “hexify” the `delete` button in `single-contact.go.html`

```html
<p>
  <button 
    hx-delete="/contacts/{{ .Data.ID }}" 
    hx-target="body"
    hx-push-url="true"
    >Delete this contact</button>
</p>
```

#### Annotations

1. `hx-delete` sends a `DELETE` request to `/contacts/5` (in our case)
2. Without the `hx-target`, the response would have replaced the `<button>`
   element that triggert the request. With `hx-target`, the right DOM element
   will be replaced
3. `hx-push-url` will see to execute the redirection even in the location bar
   will be complete and show the right URL.

But we had to update the handler for `DELETE /contacts/:id`, too:

```go
// ...
return fiberflash.WithSuccess(c, fiber.Map{
	"Status": "success",
	"Msg":    fmt.Sprintf("Successfully deleted Contact '%s' from database.", c.Params("id")),
}).Redirect("/contacts", fiber.StatusSeeOther)
// ...
```

The `fiber.StatusSeeOther` supplement was necessary to force the browser to
make a `GET` request when performing the redirection. Without it, the browser
would have continued with the latest HTTP method it was processed, with
happened to be `DELETE`.

The result:

> We [now] have a button that, all by itself, is able to issue a properly
> formatted HTTP DELETE request to the correct URL, and the UI and location bar
> are all updated correctly. This was accomplished with three declarative
> attributes placed directly on the button: hx-delete, hx-target and
> hx-push-url.

## 2024-02-23 05:36

- [x] Logger installieren
- [x] Logger mit neuem Endpoint ausprobieren; rumspielen, experimentieren.

Installed it. Logs everything now. And I found something out. `fmt.Printf()`
can be used for any special logging.
Hab ihn installiert; er loggt jetzt alles. Und ich hab was herausgefunden:

Working with HTMX will be less tedious, but it will be „different“. My biggest
problem will be thinking in variable frames. To give an example: I make a request. 
This request will result in a snippet for HTMX and shall even fire up a
flash/toast message.

Using HTMX, the flash container must be included into the snippet. CSS would define
it to be `fixed` and flashing up on the top right and vanishing again after a couple of seconds.
Yes, I can write that CSS. But will it still work when HTMX loads and
integrates the snippet? The docs say, “No problem!”. We simply have to
experiment. Therefore I made the following plan:

1. define CSS classes forl flash messages, with positioning and transitions
0. CSS-Klassen für Flash-Nachrichten einbauen
0. convert `DELETE /contacts/:id` to HTMX technology
0. `hx-swap='#contact-table'`
0. pray everythign works out well.

__[ UPDATE 12:23:]__ Im the [GitHub repository for 
Animate.css](https://github.com/animate-css/animate.css/tree/main/source)
we can steal good transitions and animations with impunity without installing
[Animate.css](https://animate.style/) itself.:w


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
