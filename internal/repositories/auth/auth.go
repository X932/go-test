package auth_repository

import (
	"database/sql"

	"go.uber.org/fx"
)

var Module = fx.Provide(NewRepo)

type Params struct {
	fx.In
	DB *sql.DB
}

type Repo interface {
	GetUserByEmail(email string) (user, error)
}

type repo struct {
	db *sql.DB
}

type user struct {
	ID       int
	Email    string
	Password string
}

func NewRepo(p Params) Repo {
	return &repo{
		db: p.DB,
	}
}

func (r *repo) GetUserByEmail(email string) (user, error) {
	row := r.db.QueryRow("select id, email, password from \"user\" where email=$1", email)

	var user user

	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return user, err
	}

	return user, nil
}
