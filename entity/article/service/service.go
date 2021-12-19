package service

import (
	"article-service/entity/article"
	"article-service/entity/article/repository"
	"time"
)

type ArticleService interface {
	PostArticle(ar *article.ArticleRequest) error
	GetArticles(ar *article.QueryParams) ([]*article.ArticleResponse, error)
}

type articleService struct {
	articleRepository repository.ArticleRepository
}

func NewArticleService(a repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: a,
	}
}

func (a *articleService) PostArticle(ar *article.ArticleRequest) error {
	ar.Created = time.Now()
	return a.articleRepository.PostArticle(ar)
}

func (a *articleService) GetArticles(ar *article.QueryParams) ([]*article.ArticleResponse, error) {
	return a.articleRepository.GetArticles(ar)
}
