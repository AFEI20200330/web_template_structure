package routes

import (
	"web_template/settings"
	"web_template/utils/logger"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine{
	 r := gin.Default()
	 r.Use(logger.GinLogger(), logger.GinRecovery(true))
	 r.GET("/version", func(c *gin.Context) {
		c.JSON(200,gin.H{
			"version": settings.Conf.Version,
		} )
	 })
	 return r
}