package article_service

import (
	"fmt"
	article_repository "test-go/internal/repositories/article"

	"go.uber.org/fx"
)

var Module = fx.Provide(NewModule)

type Params struct {
	fx.In
	ArticleRepo article_repository.Repo
}

type service struct {
	articleRepo article_repository.Repo
}

type Service interface {
	Create(article CreateParam) error
	GetArticles() []article_repository.Article
}

func NewModule(p Params) Service {
	return &service{articleRepo: p.ArticleRepo}
}

type CreateParam struct {
	Title   string
	Content string
	Tags    []string
	OwnerId int
}

func (s *service) Create(article CreateParam) error {
	repoErr := s.articleRepo.Create(article_repository.CreateParam{
		Content: article.Content,
		Title:   article.Title,
		Tags:    article.Tags,
		OwnerId: article.OwnerId,
	})

	if repoErr != nil {
		return repoErr
	}

	return nil
}

func (s *service) GetArticles() []article_repository.Article {
	articles, err := s.articleRepo.GetArticles()

	if err != nil {
		fmt.Println("==== repo error", err)
		return nil
	}

	return articles
}
