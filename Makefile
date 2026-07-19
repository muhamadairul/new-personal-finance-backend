.PHONY: run build test migrate seed

run:
	go run cmd/main.go

build:
	go build -o bin/api cmd/main.go

migrate:
	go run cli/main.go migrate

seed:
	go run cli/main.go seed
