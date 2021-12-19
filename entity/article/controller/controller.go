package controller

import (
	"article-service/entity"
	"article-service/entity/article"
	"article-service/entity/article/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	ArticleService service.ArticleService
}

func NewArticleController(e *echo.Echo, s service.ArticleService) {
	handler := &ArticleController{
		ArticleService: s,
	}
	e.POST("/articles", handler.PostArticle)
	e.GET("/articles", handler.GetArticles)
}

func (a *ArticleController) PostArticle(c echo.Context) error {
	payload := new(article.ArticleRequest)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, entity.BadRequestResponse())
	}
	if err := validator.New().Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, entity.BadRequestResponse())
	}

	err := a.ArticleService.PostArticle(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.InternalServerErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, entity.NewSuccessResponse("Post Article Success"))
}

func (a *ArticleController) GetArticles(c echo.Context) error {
	queryParams := new(article.QueryParams)
	queryParams.Author = c.QueryParam("author")
	queryParams.Query = c.QueryParam("query")
	articles, err := a.ArticleService.GetArticles(queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.InternalServerErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, entity.NewSuccessResponse(articles))
}
