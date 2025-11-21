package main

import (
	"fmt"
	database "test-go/internal/db"
	"test-go/internal/http/handlers"
	"test-go/internal/http/middlewares"
	app_http "test-go/internal/http/routes"
	"test-go/internal/repositories"
	"test-go/internal/services"
	"test-go/pkg/config"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Warning: No .env file found. Relying on OS environment.")
	}

	fx.New(
		app_http.Module,
		handlers.Module,
		services.Module,
		repositories.Module,
		database.Module,
		middlewares.Module,
		config.Module,
	).Run()
}
