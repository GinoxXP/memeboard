package api

import (
	"memeboard/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetImages(c *gin.Context) {
	images, err := service.GetAllImages()
	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Memeboard",
		"images": images,
	})
}
