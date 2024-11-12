BINPATH = bin/disapp

.PHONY: build
build: build-templ build-app

.PHONY: build-app
build-app:
	go build -o $(BINPATH) cmd/disapp/main.go

.PHONY: build-templ
build-templ:
	templ generate internal/view

.PHONY: run
run: build
	$(BINPATH)

.PHONY: watch
watch:
	$(MAKE) -j2 watch-app watch-templ

.PHONY: watch-app
watch-app:
	go run github.com/air-verse/air@latest \
	--build.cmd "$(MAKE) build-app" \
	--build.bin "$(BINPATH)" \
	--build.include_ext "go" \
	--build.exclude_dir "bin"

.PHONY: watch-templ
watch-templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false

.PHONY: clean
clean:
	rm -rf $(BINPATH)