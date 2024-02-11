ROOT = /home/matthiaslangbart/Documents/GitHub/htmx_on_fiber

build_views:
	mkdir -p views/layouts
	mkdir -p views/pages
	mkdir -p views/partials
	mkdir -p views/snippets

go_run:
	go run src/*.go
