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
	GetArticles() ([]Article, error)
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

type owner struct {
	ID        int
	FirstName string
	LastName  string
}

type Article struct {
	ID      int
	Title   string
	Content string
	Tags    []byte
	Owner   owner
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

func (r *repo) GetArticles() ([]Article, error) {
	rows, err := r.db.Query(`
		select
			a.id as article_id,
			a.title as article_title,
			a.content as article_content,
			a.tags as article_tags,
			o.id as user_id,
			o.first_name as user_first_name,
			o.last_name as user_last_name
		from article a
		left join "user" o on o.id = a.owner_id;
	`)

	var articles []Article

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var article Article
		var owner owner

		if err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.Tags,
			&owner.ID, &owner.FirstName, &owner.LastName); err != nil {
			return nil, err
		}

		article.Owner = owner
		articles = append(articles, article)
	}

	return articles, nil
}
