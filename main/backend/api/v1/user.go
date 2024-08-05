package v1

import (
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/middleware"
	"github.com/gin-gonic/gin"
)

// UserHandler api/v1/user
func UserHandler(rg *gin.RouterGroup) {
	// rg.POST() to create a new user is done in register.go

	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.VerifyPermissions(
			permissions.UserViewAll,
		),
		func(c *gin.Context) {
			// returns all users
			// TODO: implement the actual function
			c.JSON(200, gin.H{"msg": "You are able to view all users"})
		})

	rg.GET("/:id", func(c *gin.Context) {

	})

	rg.PATCH("/:id", func(c *gin.Context) {

	})

	rg.DELETE("/:id", func(c *gin.Context) {

	})
}
