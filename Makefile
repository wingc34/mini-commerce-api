.PHONY: run migrate test build

run:
	go run cmd/server/main.go

migrate:
	docker exec -i mini-commerce-api-postgres-1 psql -U postgres -d mini_commerce < migrations/001_init.sql

seed:
	docker exec -i mini-commerce-api-postgres-1 psql -U postgres -d mini_commerce < migrations/seed.sql

test:
	go test ./...

check:
	go build ./...

build:
	go build -o bin/server cmd/server/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down