package main

import (
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	storage Storage
}

func NewServer() *Server {
	return &Server{
		storage: *NewStorage(),
	}
}

func (server Server) ErrorPage(c *gin.Context, err error) {
	log.Print(err.Error())
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
		"title": "Some shit happened",
		"error": err.Error(),
	})
	c.Abort()
}

func (server Server) getAllImages(c *gin.Context) {
	images, err := server.storage.GetAllImages()
	if err != nil {
		server.ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Memeboard",
		"images": images,
	})
}

func (server Server) uploadImage(c *gin.Context) {
	file, err := c.FormFile("filename")
	if err != nil {
		server.ErrorPage(c, err)
		return
	}

	imageName := uuid.New().String() + ".png"

	imagePath := path.Join(GetImagesPath(), imageName)

	err = c.SaveUploadedFile(file, imagePath)
	if err != nil {
		server.ErrorPage(c, err)
		return
	}

	server.storage.UploadImage(imagePath)

	c.Redirect(http.StatusMovedPermanently, "/")
	c.Abort()
}
