package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestHandler(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
}
