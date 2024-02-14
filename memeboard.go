package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type DbSettings struct {
	Username string `json:username`
	Password string `json:password`
	Address  string `json:address`
	Port     int    `json:port`
	DbName   string `json:dbname`
}

func main() {

	configureDB()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/storage", "./storage")
	router.Static("/css", "./css")

	server := Server{}

	router.GET("/", server.getAllImages)
	router.POST("/upload", server.uploadImage)

	router.Run()
}

func configureDB() {
	jsonFile, err := os.Open("./dbSettings.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var settings DbSettings

	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		panic(err)
	}

	bdUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		settings.Username,
		settings.Password,
		settings.Address,
		settings.Port,
		settings.DbName)

	conn, err := pgx.Connect(context.Background(), bdUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
}
