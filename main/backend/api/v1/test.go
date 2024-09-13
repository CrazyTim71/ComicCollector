package v1

import (
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TestHandler api/v1/test
func TestHandler(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	rg.GET("admin",
		middleware.CheckJwtToken(),
		middleware.VerifyUserGroup(groups.Administrator),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello, Admin!",
			})
		})
}
