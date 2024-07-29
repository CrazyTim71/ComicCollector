package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiHandler(rg *gin.RouterGroup) {
	test := rg.Group("/test")
	TestHandler(test)

	ping := rg.Group("/ping")
	ping.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})
}
