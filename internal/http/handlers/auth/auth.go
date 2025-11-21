package auth_handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	custom_errors "test-go/internal/errors"
	"test-go/internal/response"
	auth_service "test-go/internal/services/auth"
	user_service "test-go/internal/services/user"
	"test-go/pkg/custom_regex"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewHandler)

type Params struct {
	fx.In
	AuthService auth_service.Service
	UserService user_service.Service
}

type Handler interface {
	SignIn(ctx *gin.Context)
	SignUp(ctx *gin.Context)
}

type handler struct {
	authService auth_service.Service
	userService user_service.Service
}

type SignUpCredentials struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignInCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(p Params) Handler {
	return &handler{
		authService: p.AuthService,
		userService: p.UserService,
	}
}

func validateSignInCredentials(credentials SignInCredentials) error {

	credentials.Email = strings.TrimSpace(credentials.Email)
	if !custom_regex.EmailRegex.MatchString(credentials.Email) {
		return fmt.Errorf("email is invalid")
	}

	credentials.Password = strings.TrimSpace(credentials.Password)
	if credentials.Password == "" {
		return fmt.Errorf("password is invalid")
	}

	return nil
}

func validateSignUpCredentials(credentials *SignUpCredentials) error {

	credentials.Email = strings.TrimSpace(credentials.Email)
	if !custom_regex.EmailRegex.MatchString(credentials.Email) {
		return fmt.Errorf("email is invalid")
	}

	credentials.Password = strings.TrimSpace(credentials.Password)
	if !custom_regex.PasswordRegex.MatchString(credentials.Password) {
		return fmt.Errorf("password is invalid")
	}

	credentials.FirstName = strings.TrimSpace(credentials.FirstName)
	if !custom_regex.NameRegex.MatchString(credentials.FirstName) {
		return fmt.Errorf("first_name is invalid")
	}

	credentials.LastName = strings.TrimSpace(credentials.LastName)
	if !custom_regex.NameRegex.MatchString(credentials.LastName) {
		return fmt.Errorf("last_name is invalid")
	}

	return nil
}

func (h *handler) SignIn(ctx *gin.Context) {
	var credentials SignInCredentials

	if err := ctx.BindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Body{Message: err.Error()})
		return
	}

	validationErr := validateSignInCredentials(credentials)

	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, response.Body{Message: validationErr.Error()})
		return
	}

	token, serviceErr := h.authService.SignIn(auth_service.SignInParams(credentials))

	if errors.Is(serviceErr, custom_errors.NotFoundError) {
		ctx.JSON(http.StatusNotFound, response.Body{Message: "User not found"})
		return
	}

	if serviceErr != nil {
		ctx.JSON(http.StatusBadRequest, response.Body{Message: serviceErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: "Success",
		Payload: token,
	})
}

func (h *handler) SignUp(ctx *gin.Context) {
	var credentials SignUpCredentials

	if err := ctx.BindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Body{Message: err.Error()})
		return
	}

	validationErr := validateSignUpCredentials(&credentials)

	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, response.Body{Message: validationErr.Error()})
		return
	}

	_, serviceErr := h.userService.CreateUser(user_service.NewUser(credentials))

	if serviceErr != nil {
		ctx.JSON(http.StatusBadRequest, response.Body{Message: serviceErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: "Success",
	})
}
