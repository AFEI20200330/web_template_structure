package routes

import (
	"web_template/utils/logger"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine{
	 r := gin.Default()
	 r.Use(logger.GinLogger(), logger.GinRecovery(true))
	 r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	 })
	 return r
}