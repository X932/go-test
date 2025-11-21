package config

import (
	"os"

	"go.uber.org/fx"
)

var Module = fx.Provide(LoadConfig)

type Config struct {
	DB_URL     string
	APP_PORT   string
	JWT_SECRET string
	ISSUER     string
}

func LoadConfig() *Config {
	return &Config{
		DB_URL:     os.Getenv("DATABASE_URL"),
		APP_PORT:   os.Getenv("APP_PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		ISSUER:     os.Getenv("ISSUER"),
	}
}
