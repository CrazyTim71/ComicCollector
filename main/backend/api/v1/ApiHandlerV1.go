package v1

import (
	"github.com/gin-gonic/gin"
)

func ApiHandler(rg *gin.RouterGroup) {
	// TODO: remove the test handler later
	test := rg.Group("/test")
	TestHandler(test)

	user := rg.Group("/user")
	UserHandler(user)

	login := rg.Group("/login")
	LoginHandler(login)

	register := rg.Group("/register")
	RegisterHandler(register)

	book := rg.Group("/book")
	BookHandler(book)

	author := rg.Group("/author")
	AuthorHandler(author)

	bookEdition := rg.Group("/book/edition") // TODO: check if this works
	BookEditionHandler(bookEdition)
}
