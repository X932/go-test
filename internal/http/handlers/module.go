package handlers

import (
	article_handler "test-go/internal/http/handlers/article"
	auth_handler "test-go/internal/http/handlers/auth"
	user_handler "test-go/internal/http/handlers/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	user_handler.Module,
	auth_handler.Module,
	article_handler.Module,
)
