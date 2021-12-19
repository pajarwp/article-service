package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
)

var databaseHost string
var databasePort string
var databaseUsername string
var databasePassword string
var databaseName string
var elasticHost string
var elasticPort string
var ElasticIndex string

func GetEnvVariable(path string) {
	err := godotenv.Load(path + ".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	databaseHost = os.Getenv("DATABASE_HOST")
	databasePort = os.Getenv("DATABASE_PORT")
	databaseUsername = os.Getenv("DATABASE_USERNAME")
	databasePassword = os.Getenv("DATABASE_PASSWORD")
	elasticHost = os.Getenv("ELASTIC_HOST")
	elasticPort = os.Getenv("ELASTIC_PORT")
	if os.Getenv("ENV") != "testing" {
		databaseName = os.Getenv("DATABASE_NAME")
		ElasticIndex = os.Getenv("ELASTIC_INDEX")
	} else {
		databaseName = os.Getenv("DATABASE_NAME_TEST")
		ElasticIndex = os.Getenv("ELASTIC_INDEX_TEST")
	}

}

func InitDatabase() *sql.DB {
	// init database
	dcn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", databaseUsername, databasePassword, databaseHost, databasePort, databaseName)
	dbConn, err := sql.Open(`mysql`, dcn)
	if err != nil {
		panic(fmt.Errorf("Fatal error database connection: %s \n", err))
	}
	return dbConn
}

func GetESClient() (*elastic.Client, error) {
	dcn := fmt.Sprintf("http://%s:%s", elasticHost, elasticPort)
	client, err := elastic.NewClient(elastic.SetURL(dcn),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		return nil, err
	}

	return client, err

}
