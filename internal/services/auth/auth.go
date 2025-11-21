package auth_service

import (
	"database/sql"
	"errors"
	"strconv"
	custom_errors "test-go/internal/errors"
	auth_repository "test-go/internal/repositories/auth"
	"test-go/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

var Module = fx.Provide(NewService)

type Params struct {
	fx.In
	AuthRepo auth_repository.Repo
	Config   *config.Config
}

type Service interface {
	SignIn(params SignInParams) (string, error)
}

type service struct {
	authRepo auth_repository.Repo
	config   *config.Config
}

type SignInParams struct {
	Email    string
	Password string
}

type AuthClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func NewService(p Params) Service {
	return &service{
		authRepo: p.AuthRepo,
		config:   p.Config,
	}
}

func generateToken(id int, config *config.Config) (string, error) {
	authClaims := AuthClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ISSUER,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Subject:   strconv.Itoa(id),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	token, signingTokenErr := jwtToken.SignedString([]byte(config.JWT_SECRET))

	if signingTokenErr != nil {
		return "", signingTokenErr
	}

	return token, nil
}

func (s *service) SignIn(params SignInParams) (string, error) {
	foundUser, repoErr := s.authRepo.GetUserByEmail(params.Email)

	if errors.Is(repoErr, sql.ErrNoRows) {
		return "", custom_errors.NotFoundError
	}

	if repoErr != nil {
		return "", repoErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(params.Password)); err != nil {
		return "", custom_errors.InvalidCredentialsError
	}

	token, tokenGenErr := generateToken(foundUser.ID, s.config)

	if tokenGenErr != nil {
		return "", tokenGenErr
	}

	return token, nil
}
