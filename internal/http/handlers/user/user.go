package user_handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	custom_errors "test-go/internal/errors"
	user_repository "test-go/internal/repositories/user"
	"test-go/internal/response"
	user_service "test-go/internal/services/user"
	"test-go/pkg/constants"
	"test-go/pkg/custom_regex"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewUserController)

const DEFAULT_ROWS_LIMIT = "10"

type Handler interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type handler struct {
	userService user_service.Service
}

type Params struct {
	fx.In
	UserService user_service.Service
}

type UserDto struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreateUserDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserParams struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func NewUserController(p Params) Handler {
	return &handler{
		userService: p.UserService,
	}
}

func validateUserParams(user UserParams) error {
	if !custom_regex.EmailRegex.MatchString(user.Email) {
		return fmt.Errorf("email is invalid")
	}

	user.FirstName = strings.TrimSpace(user.FirstName)
	if !custom_regex.NameRegex.MatchString(user.FirstName) {
		return fmt.Errorf("first_name is invalid")
	}

	user.LastName = strings.TrimSpace(user.LastName)
	if !custom_regex.NameRegex.MatchString(user.LastName) {
		return fmt.Errorf("last_name is invalid")
	}

	return nil
}

func (h *handler) CreateUser(c *gin.Context) {
	var newUser CreateUserDto

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	validationErrMessage := validateUserParams(UserParams(newUser))

	if validationErrMessage != nil {
		c.JSON(http.StatusBadRequest, response.Body{Message: validationErrMessage.Error()})
		return
	}

	rowsAffected, createUserErr := h.userService.CreateUser(user_service.NewUser{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
	})

	if createUserErr != nil {
		c.JSON(http.StatusBadRequest, response.Body{
			Message: createUserErr.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, rowsAffected)
}

func (h *handler) GetUsers(c *gin.Context) {
	paramRowsLimit := c.Query("rows_limit")
	paramRowsOffset := c.Query("rows_offset")

	if strings.TrimSpace(paramRowsLimit) == "" {
		paramRowsLimit = DEFAULT_ROWS_LIMIT
	}

	if strings.TrimSpace(paramRowsOffset) == "" {
		paramRowsOffset = "0"
	}

	rowsLimit, rowsLimitErr := strconv.Atoi(paramRowsLimit)
	rowsOffset, rowsOffsetErr := strconv.Atoi(paramRowsOffset)

	if rowsLimitErr != nil || rowsOffsetErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Body{Message: "Invalid params"})
		return
	}

	res := h.userService.GetUsers(rowsLimit, rowsOffset)

	var users []UserDto
	for _, u := range res {
		users = append(users, UserDto(u))
	}

	c.JSON(http.StatusOK, response.Body{Message: "Success", Payload: users})
}

func (h *handler) GetUser(c *gin.Context) {
	paramID := c.Param("id")

	userID, _ := c.Get(constants.USER_CTX_KEY)

	fmt.Printf("====== userId has value=%v, type=%T\n", userID, userID)

	if paramID == "" {
		c.JSON(http.StatusBadRequest, "user ID is required")
		return
	}

	ID, strConvErr := strconv.Atoi(paramID)
	if strConvErr != nil {
		c.JSON(http.StatusBadRequest, "user ID is invalid")
		return
	}

	user, serviceErr := h.userService.GetUser(ID)

	if serviceErr != nil {
		c.JSON(http.StatusBadRequest, response.Body{Message: serviceErr.Error()})
		return
	}

	c.JSON(http.StatusOK, UserDto(user))
}

func (h *handler) DeleteUser(c *gin.Context) {
	paramID := c.Param("id")
	id, strconvErr := strconv.Atoi(paramID)

	if strconvErr != nil {
		c.JSON(http.StatusBadRequest, response.Body{
			Message: "Invalid user ID",
		})
		return
	}

	execErr := h.userService.DeleteUser(id)

	if execErr != nil {
		c.JSON(http.StatusNotFound, response.Body{
			Message: execErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Body{
		Message: "Success",
	})
}

func (h *handler) UpdateUser(c *gin.Context) {
	var editedUserParams UserDto

	if err := c.BindJSON(&editedUserParams); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	validationErrMessage := validateUserParams(UserParams{
		FirstName: editedUserParams.FirstName,
		LastName:  editedUserParams.LastName,
		Email:     editedUserParams.Email,
	})

	if validationErrMessage != nil {
		c.JSON(http.StatusBadRequest, response.Body{
			Message: validationErrMessage.Error(),
		})
		return
	}

	updatedUser, serviceErr := h.userService.UpdateUser(user_repository.User(editedUserParams))

	if errors.Is(serviceErr, custom_errors.NotFoundError) {
		c.JSON(http.StatusNotFound, response.Body{Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserDto(updatedUser))
}
