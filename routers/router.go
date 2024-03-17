package routers

import (
	"memeboard/routers/api"
	"memeboard/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	service.Init()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/storage", "./storage")
	router.Static("/css", "./css")

	router.GET("/", api.GetImages)
	router.GET("/image", api.GetImage)
	router.POST("/upload", api.UploadImage)

	router.Run()
}
