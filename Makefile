include .env

build:
	@go build -o bin/career-path cmd/app/main.go

run: build
	@./bin/career-path

test:
	@go test ./... -v

migrate-create:
	migrate create -ext sql -dir database/migrations $(name)

migrate-up:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path database/migrations up

migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -path database/migrations down