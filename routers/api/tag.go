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

	service.AddTagToImage(tag, imageID)

	c.Redirect(http.StatusMovedPermanently, "/image?id="+c.PostForm("imageID"))
	c.Abort()
}
