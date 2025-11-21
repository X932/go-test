package repositories

import (
	auth_repository "test-go/internal/repositories/auth"
	user_repository "test-go/internal/repositories/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	user_repository.Module,
	auth_repository.Module,
)
