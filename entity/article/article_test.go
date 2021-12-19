package article_test

import (
	"article-service/config"
	"article-service/entity"
	"article-service/entity/article"
	"article-service/entity/article/controller"
	"article-service/entity/article/repository"
	"article-service/entity/article/service"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

var (
	articleRequest = &article.ArticleRequest{
		Author: "John Doe",
		Title:  "Crypto",
		Body:   "Buy Bitcoin",
	}

	articleRequest2 = &article.ArticleRequest{
		Author: "Jane Doe",
		Title:  "Economy",
		Body:   "Buy crypto",
	}

	articleRequest3 = &article.ArticleRequest{
		Author: "John Doe",
		Title:  "Animal",
		Body:   "Tiger",
	}

	articleRequestBadRequest = &article.ArticleRequest{
		Title: "Title Test 1",
		Body:  "Body Test 1",
	}

	handler  controller.ArticleController
	e        *echo.Echo
	esClient *elastic.Client
)

func TestMain(m *testing.M) {

	c := exec.Command("bash", "-c", "make -C ../../ migrateuptest")
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	dirname, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Open(path.Join(dirname, "../../"))
	config.GetEnvVariable(dir.Name() + "/")
	dbConn := config.InitDatabase()
	esClient, err = config.GetESClient()
	if err != nil {
		log.Fatal("Error initializing elastic : ", err)
	}
	e = echo.New()

	repo := repository.NewArticleRepository(dbConn, esClient)
	service := service.NewArticleService(repo)
	handler = controller.ArticleController{
		ArticleService: service,
	}
	exitVal := m.Run()
	c = exec.Command("bash", "-c", "make -C ../../ migratedowntest")
	c.Stdin = strings.NewReader("y")
	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	_, err = esClient.DeleteIndex(config.ElasticIndex).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitVal)
}

func TestPostArticle(t *testing.T) {
	bodyRequest := articleRequest
	jsonBody, _ := json.Marshal(bodyRequest)
	req, _ := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	_ = handler.PostArticle(c)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestPostArticle2(t *testing.T) {
	bodyRequest := articleRequest2
	jsonBody, _ := json.Marshal(bodyRequest)
	req, _ := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	_ = handler.PostArticle(c)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestPostArticle3(t *testing.T) {
	bodyRequest := articleRequest3
	jsonBody, _ := json.Marshal(bodyRequest)
	req, _ := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	_ = handler.PostArticle(c)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestPostArticleBadRequest(t *testing.T) {
	bodyRequest := articleRequestBadRequest
	jsonBody, _ := json.Marshal(bodyRequest)
	req, _ := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	_ = handler.PostArticle(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

}

func TestPostArticleNoContentType(t *testing.T) {
	bodyRequest := articleRequest
	jsonBody, _ := json.Marshal(bodyRequest)
	req, _ := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(jsonBody)))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	_ = handler.PostArticle(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

}

func TestGetAllArticle(t *testing.T) {
	ctx := context.Background()
	esClient.Refresh(config.ElasticIndex).Do(ctx)
	req, _ := http.NewRequest(echo.GET, "/articles", strings.NewReader(""))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	_ = handler.GetArticles(c)
	bodyByte := rec.Body.Bytes()
	var body entity.DefaultResponse
	_ = json.Unmarshal(bodyByte, &body)
	byteData, _ := json.Marshal(body.Data)
	var data []*article.ArticleResponse
	_ = json.Unmarshal(byteData, &data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, "Animal", data[0].Title)

}

func TestGetArticleByAuthor(t *testing.T) {
	ctx := context.Background()
	esClient.Refresh(config.ElasticIndex).Do(ctx)
	req, _ := http.NewRequest(echo.GET, "/articles", strings.NewReader(""))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	c.QueryParams().Add("author", "John Doe")
	_ = handler.GetArticles(c)
	bodyByte := rec.Body.Bytes()
	var body entity.DefaultResponse
	_ = json.Unmarshal(bodyByte, &body)
	byteData, _ := json.Marshal(body.Data)
	var data []*article.ArticleResponse
	_ = json.Unmarshal(byteData, &data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, "Animal", data[0].Title)

}

func TestGetArticleByQuery(t *testing.T) {
	ctx := context.Background()
	esClient.Refresh(config.ElasticIndex).Do(ctx)
	req, _ := http.NewRequest(echo.GET, "/articles", strings.NewReader(""))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	c.QueryParams().Add("query", "Crypto")
	_ = handler.GetArticles(c)
	bodyByte := rec.Body.Bytes()
	var body entity.DefaultResponse
	_ = json.Unmarshal(bodyByte, &body)
	byteData, _ := json.Marshal(body.Data)
	var data []*article.ArticleResponse
	_ = json.Unmarshal(byteData, &data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, "Economy", data[0].Title)

}

func TestGetArticleByQueryAndAuthor(t *testing.T) {
	ctx := context.Background()
	esClient.Refresh(config.ElasticIndex).Do(ctx)
	req, _ := http.NewRequest(echo.GET, "/articles", strings.NewReader(""))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles")
	c.QueryParams().Add("query", "Crypto")
	c.QueryParams().Add("author", "John Doe")
	_ = handler.GetArticles(c)
	bodyByte := rec.Body.Bytes()
	var body entity.DefaultResponse
	_ = json.Unmarshal(bodyByte, &body)
	byteData, _ := json.Marshal(body.Data)
	var data []*article.ArticleResponse
	_ = json.Unmarshal(byteData, &data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "Crypto", data[0].Title)

}
