package v1

import (
	"github.com/gin-gonic/gin"
)

func LoginHandler(rg *gin.RouterGroup) {
	rg.POST("", func(c *gin.Context) {
		// TODO: generate jwt
		// TODO: make middleware to check the jwt token on every endpoint
	})

}
