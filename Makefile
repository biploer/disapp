.PHONY: build
build: build-templ build-css build-app

.PHONY: build-app
build-app:
	go build -o ./bin/disapp cmd/disapp/main.go

.PHONY: build-templ
build-templ:
	templ generate

.PHONY: build-css
build-css:
	npm --prefix web run build:css

.PHONY: run
run: build
	./bin/disapp