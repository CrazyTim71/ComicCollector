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

	author := rg.Group("/book/author")
	AuthorHandler(author)

	bookEdition := rg.Group("/book/edition")
	BookEditionHandler(bookEdition)

	bookType := rg.Group("/book/type")
	BookTypeHandler(bookType)

	location := rg.Group("/book/location")
	LocationHandler(location)

	owner := rg.Group("/book/owner")
	OwnerHandler(owner)

	publisher := rg.Group("/book/publisher")
	PublisherHandler(publisher)

	// admin only
	permissions := rg.Group("/permission")
	PermissionsHandler(permissions)

	// admin only
	role := rg.Group("/role")
	RoleHandler(role)

}
