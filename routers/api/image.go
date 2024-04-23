package api

import (
	"memeboard/service"
	"net/http"
	"strconv"
	"strings"

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

func GetImagesForTags(c *gin.Context) {
	search := c.PostForm("search")
	tagnames := strings.Split(search, ", ")
	images, err := service.GetImagesByTags(tagnames)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Memeboard",
		"images": images,
	})
	c.Abort()
}

func GetImage(c *gin.Context) {
	query := c.Request.URL.Query()
	id, err := strconv.Atoi(query.Get("id"))

	if err != nil {
		ErrorPage(c, err)
		return
	}

	image, err := service.GetImage(id)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	tags, err := service.GetAttachedTags(image.ID)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "image.tmpl", gin.H{
		"title": "Image",
		"image": image,
		"tags":  tags,
	})
}
