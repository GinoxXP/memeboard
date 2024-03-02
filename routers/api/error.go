package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorPage(c *gin.Context, err error) {
	log.Print(err.Error())
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
		"title": "Some shit happened",
		"error": err.Error(),
	})
	c.Abort()
}
