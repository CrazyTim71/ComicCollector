package v1

import "github.com/gin-gonic/gin"

// UserHandler api/v1/user
func UserHandler(rg *gin.RouterGroup) {
	// rg.POST() to create a new user is done in register.go
	
	rg.GET("/:id", func(c *gin.Context) {

	})

	rg.PATCH("/:id", func(c *gin.Context) {

	})

	rg.DELETE("/:id", func(c *gin.Context) {

	})
}
