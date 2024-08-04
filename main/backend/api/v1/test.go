package v1

import (
	"ComicCollector/main/backend/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestHandler(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	rg.GET("admin", middleware.CheckJwtToken(), middleware.VerifyAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Admin!",
		})
	})
}
