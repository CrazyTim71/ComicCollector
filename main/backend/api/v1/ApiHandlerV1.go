package v1

import (
	"github.com/gin-gonic/gin"
)

func ApiHandler(rg *gin.RouterGroup) {
	test := rg.Group("/test")
	TestHandler(test)

	user := rg.Group("/user")
	UserHandler(user)

	login := rg.Group("/login")
	LoginHandler(login)

	register := rg.Group("/register")
	RegisterHandler(register)
}
