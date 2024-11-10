.PHONY: build
build:
	go build -o ./bin/disapp cmd/disapp/main.go

.PHONY: run
run: build
	./bin/disapp
