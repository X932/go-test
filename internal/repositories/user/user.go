package user_repository

import (
	"database/sql"
	"errors"
	"fmt"
	custom_errors "test-go/internal/errors"

	"go.uber.org/fx"
)

var Module = fx.Provide(NewModule)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

type CreateUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  []byte
}

type Repo interface {
	CreateUser(newUser CreateUser) (int64, error)
	GetUsers(rowsLimit, rowsOffset int) ([]User, error)
	GetUser(id int) (User, error)
	DeleteUser(id int) (int64, error)
	UpdateUser(user User) (User, error)
}

type repo struct {
	db *sql.DB
}

type Params struct {
	fx.In
	DB *sql.DB
}

func NewModule(p Params) Repo {
	return &repo{
		db: p.DB,
	}
}

func (r *repo) CreateUser(newUser CreateUser) (int64, error) {
	sqlResult, err := r.db.Exec("insert into \"user\" (email, first_name, last_name, password) values ($1, $2, $3, $4)",
		newUser.Email, newUser.FirstName, newUser.LastName, newUser.Password)

	if err != nil {
		fmt.Printf("err = %+v", err)
		return 0, err
	}

	rowsAffected, rowsAffectedErr := sqlResult.RowsAffected()

	if rowsAffectedErr != nil {
		fmt.Printf("rowsAffectedErr = %+v", rowsAffectedErr)
		return 0, rowsAffectedErr
	}

	if rowsAffected == 0 {
		fmt.Printf("rows Affected -- %v", rowsAffected)
		return 0, rowsAffectedErr
	}

	return rowsAffected, nil
}

func (r *repo) GetUsers(rowsLimit, rowsOffset int) ([]User, error) {
	fmt.Printf("====== %v - %v\n", rowsLimit, rowsOffset)
	rows, err := r.db.Query(`
		select 
			id,
			first_name,
			last_name,
			email
		from "user" 
		order by id
		limit $1 offset $2;`,
		rowsLimit, rowsOffset*rowsLimit)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *repo) GetUser(id int) (User, error) {
	row := r.db.QueryRow(`
		select 
			id,
			first_name,
			last_name,
			email
		from "user" 
		where id=$1;`,
		id)

	var foundUser User

	if err := row.Scan(&foundUser.ID, &foundUser.FirstName, &foundUser.LastName, &foundUser.Email); err != nil {
		return foundUser, fmt.Errorf("%w", err)
	}

	return foundUser, nil
}

func (r *repo) UpdateUser(user User) (User, error) {
	var updatedUser User

	tx, txErr := r.db.Begin()

	if txErr != nil {
		return updatedUser, txErr
	}

	defer tx.Rollback()

	sqlResult, sqlErr := tx.Exec(`
	update "user" 
	set 
		first_name=$1, 
		last_name=$2, 
		email=$3 
	where 
		id=$4`,
		user.FirstName, user.LastName, user.Email, user.ID)

	if sqlErr != nil {
		return updatedUser, sqlErr
	}

	rowsAffected, rowsAffectedErr := sqlResult.RowsAffected()

	if rowsAffectedErr != nil {
		return updatedUser, rowsAffectedErr
	}

	if rowsAffected == 0 {
		return updatedUser, custom_errors.NotFoundError
	}

	row := tx.QueryRow("select * from \"user\" where id=$1", user.ID)

	if err := row.Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
	); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return updatedUser, custom_errors.NotFoundError
		}

		return updatedUser, err
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		return updatedUser, commitErr
	}

	return updatedUser, nil
}

func (r *repo) DeleteUser(id int) (int64, error) {
	sqlResult, sqlErr := r.db.Exec("delete from \"user\" where id=$1", id)

	if sqlErr != nil {
		return 0, sqlErr
	}

	rowsAffected, rowsAffectedErr := sqlResult.RowsAffected()

	return rowsAffected, rowsAffectedErr
}
