package services

import (
	article_service "test-go/internal/services/article"
	auth_service "test-go/internal/services/auth"
	user_service "test-go/internal/services/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	user_service.Module,
	auth_service.Module,
	article_service.Module,
)
