package database

import (
	"context"
	"database/sql"
	"fmt"
	"test-go/pkg/config"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewConnection)

type Params struct {
	fx.In
	fx.Lifecycle

	Config *config.Config
}

func NewConnection(p Params) (*sql.DB, error) {
	connStr := p.Config.DB_URL
	db, openConnErr := sql.Open("postgres", connStr)

	if openConnErr != nil {
		return nil, openConnErr
	}

	p.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			pingErr := db.Ping()
			if pingErr != nil {
				fmt.Println("Cannot connect to DB")
				return pingErr
			}

			fmt.Println("DB connected !")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Closing DB connection !")
			return db.Close()
		},
	})

	return db, nil
}
