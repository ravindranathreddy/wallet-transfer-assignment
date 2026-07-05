.PHONY: build run test fmt lint tidy \
	docker-up docker-down docker-logs \
	migrate-create migrate-up migrate-down

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	gofmt -l .

lint:
	golangci-lint run

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Usage: make migrate-create NAME=create_wallets_table
migrate-create:
	docker run --rm -v $(CURDIR)/migrations:/migrations migrate/migrate:v4.17.1 \
		create -ext sql -dir /migrations -seq $(NAME)

migrate-up:
	docker compose --profile tools run --rm migrate up

migrate-down:
	docker compose --profile tools run --rm migrate down 1
