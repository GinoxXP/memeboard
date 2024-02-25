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

func (server *Server) ErrorPage(c *gin.Context, err error) {
	log.Print(err.Error())
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
		"title": "Some shit happened",
		"error": err.Error(),
	})
	c.Abort()
}

func (server *Server) GetGalery(c *gin.Context) {
	thumbnails, err := server.storage.GetAllThumbnails()
	if err != nil {
		server.ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Memeboard",
		"images": thumbnails,
	})
}

func (server *Server) UploadImage(c *gin.Context) {
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

	err = server.storage.UploadImage(imagePath)
	if err != nil {
		server.ErrorPage(c, err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
	c.Abort()
}
