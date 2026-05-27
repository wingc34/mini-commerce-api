.PHONY: run migrate test build

run:
	go run cmd/server/main.go

migrate:
	docker exec -i mini-commerce-api-postgres-1 psql -U postgres -d mini_commerce < migrations/001_init.sql

test:
	go test ./...

build:
	go build -o bin/server cmd/server/main.go