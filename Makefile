BINPATH = bin/burning-notes

.PHONY: build
build: build-templ build-css build-app

.PHONY: build-app
build-app:
	go build -o $(BINPATH) cmd/burning-notes/main.go

.PHONY: build-templ
build-templ:
	go get -u github.com/a-h/templ
	go run github.com/a-h/templ/cmd/templ@latest generate

.PHONY: build-css
build-css:
	npm --prefix web run build:css

.PHONY: run
run: build
	$(BINPATH)

.PHONY: watch
watch:
	$(MAKE) -j3 watch-css watch-templ watch-app

.PHONY: watch-app
watch-app:
	go run github.com/air-verse/air@latest \
	--build.cmd "$(MAKE) build-app" \
	--build.bin $(BINPATH) \
	--build.include_ext "go" \
	--build.exclude_dir "web/node_modules, bin, web"

.PHONY: watch-templ
watch-templ:
	go get -u github.com/a-h/templ
	go run github.com/a-h/templ/cmd/templ@latest generate \
	--watch \
	--proxy="http://localhost:8080" \
	--open-browser=false

.PHONY: watch-css
watch-css:
	npm --prefix web run watch:css
