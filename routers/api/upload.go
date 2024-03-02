package api

import (
	"memeboard/service"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("filename")
	if err != nil {
		ErrorPage(c, err)
		return
	}

	imageName := uuid.New().String() + ".png"

	imagePath := path.Join(service.GetImagesPath(), imageName)

	err = c.SaveUploadedFile(file, imagePath)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	err = service.UploadImage(imagePath)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
	c.Abort()
}
