package main

import (
	"github.com/gin-gonic/gin"
)

type DbSettings struct {
	Username string `json:username`
	Password string `json:password`
	Address  string `json:address`
	Port     int    `json:port`
	DbName   string `json:dbname`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/storage", "./storage")
	router.Static("/css", "./css")

	server := NewServer()

	router.GET("/", server.getAllImages)
	router.POST("/upload", server.uploadImage)

	router.Run()
}
