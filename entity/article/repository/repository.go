package repository

import (
	"article-service/entity/article"
	"context"
	"database/sql"
	"encoding/json"

	"article-service/config"

	"github.com/olivere/elastic/v7"
)

type ArticleRepository interface {
	PostArticle(ar *article.ArticleRequest) error
	GetArticles(ar *article.QueryParams) ([]*article.ArticleResponse, error)
}

type articleRepository struct {
	db *sql.DB
	es *elastic.Client
}

func NewArticleRepository(db *sql.DB, es *elastic.Client) ArticleRepository {
	return &articleRepository{
		db: db,
		es: es,
	}
}

func (a *articleRepository) PostArticle(ar *article.ArticleRequest) error {
	tx, err := a.db.Begin()
	query := "INSERT articles SET author=?, title=?, body=?, created=?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ar.Author, ar.Title, ar.Body, ar.Created)
	if err != nil {
		return err
	}
	err = a.insertToElastic(ar)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil

}

func (a *articleRepository) insertToElastic(ar *article.ArticleRequest) error {
	ctx := context.Background()
	dataJSON, err := json.Marshal(ar)
	js := string(dataJSON)
	_, err = a.es.Index().
		Index(config.ElasticIndex).
		BodyJson(js).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil

}

func (a *articleRepository) GetArticles(ar *article.QueryParams) ([]*article.ArticleResponse, error) {
	articles := make([]*article.ArticleResponse, 0)
	ctx := context.Background()
	boolQuery := elastic.NewBoolQuery()
	mustQueries := make([]elastic.Query, 0)
	if ar.Author != "" {
		mustQueries = append(mustQueries, elastic.NewTermQuery("author.keyword", ar.Author))
	}
	if ar.Query != "" {
		mustQueries = append(mustQueries, elastic.NewMultiMatchQuery(ar.Query, "title", "body").Operator("and"))
	}
	boolQuery.Must(mustQueries...)
	searchService := a.es.Search().Index(config.ElasticIndex).Query(boolQuery).Sort("Created", false)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return articles, err
	}

	for _, hit := range searchResult.Hits.Hits {
		article := new(article.ArticleResponse)
		err := json.Unmarshal(hit.Source, article)
		if err != nil {
			return articles, err
		}

		articles = append(articles, article)
	}

	return articles, nil

}
