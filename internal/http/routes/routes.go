package app_http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	auth_handler "test-go/internal/http/handlers/auth"
	user_handler "test-go/internal/http/handlers/user"
	"test-go/internal/http/middlewares"
	"test-go/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Invoke(NewHttpServer)

type Params struct {
	fx.In
	fx.Lifecycle
	UserHandler   user_handler.Handler
	AuthHandler   auth_handler.Handler
	MiddlewareSet middlewares.MiddlewareSet
	Config        *config.Config
}

func NewHttpServer(p Params) *http.Server {
	engine := gin.Default()

	engine.Use(p.MiddlewareSet.Logging())

	engine.POST("/sign-in", p.AuthHandler.SignIn)
	engine.POST("/sign-up", p.AuthHandler.SignUp)

	protected := engine.Group("/")
	protected.Use(p.MiddlewareSet.AuthMiddleware())
	{
		protected.GET("/users", p.UserHandler.GetUsers)
		protected.GET("/users/:id", p.UserHandler.GetUser)
		protected.POST("/users", p.UserHandler.CreateUser)
		protected.PUT("/users", p.UserHandler.UpdateUser)
		protected.DELETE("/users/:id", p.UserHandler.DeleteUser)
	}

	srv := &http.Server{
		Addr:    p.Config.APP_PORT,
		Handler: engine,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)

			if err != nil {
				return err
			}

			fmt.Printf("Server started at %v\n", srv.Addr)

			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("========== stop http-server ==============")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
