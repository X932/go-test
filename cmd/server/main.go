package main

import (
	database "test-go/internal/db"
	"test-go/internal/http/handlers"
	"test-go/internal/http/middlewares"
	app_http "test-go/internal/http/routes"
	"test-go/internal/repositories"
	"test-go/internal/services"
	"test-go/pkg/config"

	"go.uber.org/fx"
)

func main() {

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
