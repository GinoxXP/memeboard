package api

import (
	"memeboard/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddTag(c *gin.Context) {
	tag := c.PostForm("tag")
	imageID, err := strconv.Atoi(c.PostForm("imageID"))

	if err != nil {
		ErrorPage(c, err)
		return
	}

	err = service.AddTagToImage(tag, imageID)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/image?id="+c.PostForm("imageID"))
	c.Abort()
}

func RemoveTag(c *gin.Context) {
	tagId, err := strconv.Atoi(c.PostForm("tagID"))
	if err != nil {
		ErrorPage(c, err)
		return
	}

	imageId, err := strconv.Atoi(c.PostForm("imageID"))
	if err != nil {
		ErrorPage(c, err)
		return
	}

	err = service.RemoveTagFromImage(tagId, imageId)
	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/image?id="+c.PostForm("imageID"))
	c.Abort()
}
