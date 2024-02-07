package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/storage", "./storage")
	r.Static("/css", "./css")

	f, err := os.Open("./storage")
	if err != nil {
		log.Print(err)
		return
	}

	files, err := f.Readdir(0)
	if err != nil {
		log.Print(err)
		return
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
	r.Run()
}
