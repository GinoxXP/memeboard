package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/storage", "./storage")
	r.Static("/css", "./css")

	f, err := os.Open("./storage")
	if err != nil {
		panic(err)
	}

	files, err := f.Readdir(0)
	if err != nil {
		panic(err)
	}

	images := make([]string, len(files))

	for i := 0; i < len(files); i++ {
		images[i] = "./storage/" + files[i].Name()
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":  "Memeboard",
			"images": images,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("filename")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			log.Print(err.Error())
			return
		}

		err = c.SaveUploadedFile(file, "./storage/"+file.Filename)
		if err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			log.Print(err.Error())
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
	})
	r.Run()
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
