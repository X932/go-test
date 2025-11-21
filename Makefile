# app
run-app:
	go run cmd/server/main.go
test:
	go test -v ./...

# DB
DB_URL=postgres://postgres:Zdarova1@localhost:5432/postgres?sslmode=disable

NAME :=

create-migration:
ifndef NAME
	$(error NAME is not set. Usage: make create-migration NAME=my_new_feature)
endif
	migrate create -ext sql -dir internal/db/migrations -seq $(NAME)

migrate:
	migrate -database $(DB_URL) -path internal/db/migrations up