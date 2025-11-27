include .env

# app
run-app:
	go run cmd/server/main.go
test:
	go test -v ./...

# DB
DB_URL ?= $(DATABASE_URL)

NAME :=

create-migration:
ifndef NAME
	$(error NAME is not set. Usage: make create-migration NAME=my_new_feature)
endif
	migrate create -ext sql -dir internal/db/migrations -seq $(NAME)

migrate:
ifndef DB_URL
	$(error DB_URL is not set: Usage make migrate DB_URL=test_db_url)
endif
	migrate -database $(DB_URL) -path internal/db/migrations up
