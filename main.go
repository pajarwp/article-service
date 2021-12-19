package main

import (
	"article-service/config"
	"article-service/entity/article/controller"
	"article-service/entity/article/repository"
	"article-service/entity/article/service"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {

	config.GetEnvVariable("")
	dbConn := config.InitDatabase()
	esClient, err := config.GetESClient()
	if err != nil {
		log.Fatal("Error initializing elastic : ", err)
	}
	e := echo.New()

	repo := repository.NewArticleRepository(dbConn, esClient)
	service := service.NewArticleService(repo)
	controller.NewArticleController(e, service)

	e.Start(":8000")
}
