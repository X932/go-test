package middlewares

import (
	authmiddleware "test-go/internal/http/middlewares/auth_middleware"
	logging_middleware "test-go/internal/http/middlewares/logging"
	user_service "test-go/internal/services/user"
	"test-go/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewMiddlewareSet)

type MiddlewareSet interface {
	Logging() gin.HandlerFunc
	AuthMiddleware() gin.HandlerFunc
}

type Params struct {
	fx.In
	UserService user_service.Service
	Config      *config.Config
}

type middlewareSet struct {
	userService user_service.Service
	JWT_SECRET  string
}

func NewMiddlewareSet(p Params) MiddlewareSet {
	return &middlewareSet{
		userService: p.UserService,
		JWT_SECRET:  p.Config.JWT_SECRET,
	}
}

func (m *middlewareSet) AuthMiddleware() gin.HandlerFunc {
	return authmiddleware.Authmiddleware(m.userService, m.JWT_SECRET)
}

func (m *middlewareSet) Logging() gin.HandlerFunc {
	return logging_middleware.Logging()
}
