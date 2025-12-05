package article_handler

import (
	"net/http"
	"strings"
	"test-go/internal/response"
	article_service "test-go/internal/services/article"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewModule)

type Params struct {
	fx.In
	ArticleService article_service.Service
}

type Handler interface {
	Create(c *gin.Context)
}

type handler struct {
	articleService article_service.Service
}

type CreateArticleDto struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	OwnerId int      `json:"owner_id"`
	Tags    []string `json:"tags"`
}

func NewModule(p Params) Handler {
	return &handler{articleService: p.ArticleService}
}

func (h *handler) Create(c *gin.Context) {
	var article CreateArticleDto

	if err := c.BindJSON(&article); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Body{Message: err.Error()})
		return
	}

	article.Title = strings.TrimSpace(article.Title)
	article.Content = strings.TrimSpace(article.Content)

	if len(article.Title) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Body{Message: "Title is required"})
		return
	}

	if len(article.Content) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Body{Message: "Content is required"})
		return
	}

	if !(article.OwnerId > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Body{Message: "Owner ID is invalid"})
		return
	}

	err := h.articleService.Create(article_service.CreateParam{
		Title:   article.Title,
		Content: article.Content,
		Tags:    article.Tags,
		OwnerId: article.OwnerId,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Body{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.Body{Message: "Success"})
}
