package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

func (server Server) getAllImages(c *gin.Context) {
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

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Memeboard",
		"images": images,
	})
}

func (server Server) uploadImage(c *gin.Context) {
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
}
