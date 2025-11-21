package user_service

import (
	"fmt"
	user_repository "test-go/internal/repositories/user"

	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

var Module = fx.Provide(NewUserService)

type Service interface {
	CreateUser(newUser NewUser) (int64, error)
	GetUsers(rowsLimit, rowsOffset int) []user_repository.User
	GetUser(id int) (user_repository.User, error)
	DeleteUser(id int) error
	UpdateUser(userUpdatedParams user_repository.User) (user_repository.User, error)
}

type Params struct {
	fx.In

	UserRepo user_repository.Repo
}

type service struct {
	userRepo user_repository.Repo
}

type NewUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func NewUserService(p Params) Service {
	return &service{
		userRepo: p.UserRepo,
	}
}

func (s *service) GetUser(id int) (user_repository.User, error) {
	user, err := s.userRepo.GetUser(id)

	if err != nil {
		return user_repository.User{}, err
	}

	return user, nil
}

func (s *service) GetUsers(rowsLimit, rowsOffset int) []user_repository.User {
	users, err := s.userRepo.GetUsers(rowsLimit, rowsOffset)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return users
}

func (s *service) CreateUser(newUser NewUser) (int64, error) {
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		return 0, hashErr
	}

	return s.userRepo.CreateUser(user_repository.CreateUser{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Password:  hashedPassword,
	})
}

func (s *service) DeleteUser(id int) error {
	rowsAffected, rowsAffectedErr := s.userRepo.DeleteUser(id)

	if rowsAffected == 0 {
		return fmt.Errorf("Not Found")
	}

	return rowsAffectedErr
}

func (s *service) UpdateUser(userUpdatedParams user_repository.User) (user_repository.User, error) {
	updatedUser, repoErr := s.userRepo.UpdateUser(userUpdatedParams)

	if repoErr != nil {
		return updatedUser, repoErr
	}

	return updatedUser, nil
}
