package article_repository

import (
	"database/sql"
	"fmt"

	"github.com/bytedance/sonic"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewModule)

type Repo interface {
	Create(params CreateParam) error
}

type Params struct {
	fx.In
	DB *sql.DB
}

type repo struct {
	db *sql.DB
}

type CreateParam struct {
	Title   string
	Content string
	Tags    []string
	OwnerId int
}

func NewModule(p Params) Repo {
	return &repo{db: p.DB}
}

func (r *repo) Create(params CreateParam) error {
	var tagsSqlValue any

	if len(params.Tags) > 0 {
		jsonTags, tagsConvertionErr := sonic.Marshal(params.Tags)

		if tagsConvertionErr != nil {
			return tagsConvertionErr
		}

		tagsSqlValue = jsonTags
	} else {
		tagsSqlValue = nil
	}

	sqlResult, sqlErr := r.db.Exec(`
		insert into article 
			(title, content, tags, owner_id) 
		values ($1, $2, $3, $4);
	`, params.Title, params.Content, tagsSqlValue, params.OwnerId)

	if sqlErr != nil {
		return sqlErr
	}

	if rowsAffected, err := sqlResult.RowsAffected(); err != nil || rowsAffected == 0 {
		return fmt.Errorf("Not affected")
	}

	return nil
}
