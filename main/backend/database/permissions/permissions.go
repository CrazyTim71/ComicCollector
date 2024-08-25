package permissions

type Permission struct {
	Name        string
	Description string
}

var (
	GlobalEnableEndpointAccess = Permission{
		Name:        "global:EnableEndpointAccess",
		Description: "Gives an user the ability to request data from the endpoints",
	}
)

var (
	UserViewSelf = Permission{
		Name:        "user:viewSelf",
		Description: "Allows a user to view their own information",
	}
	UserViewAll = Permission{
		Name:        "user:viewAll",
		Description: "Allows viewing all user information (admin only)",
	}

	UserModifySelf = Permission{
		Name:        "user:modifySelf",
		Description: "Allows a user to modify their own user information",
	}
	UserModifyAll = Permission{
		Name:        "user:modifyAll",
		Description: "Allows modifying all user information (admin only)",
	}

	UserDeleteSelf = Permission{
		Name:        "user:deleteSelf",
		Description: "Allows a user to delete their own account",
	}
	UserDeleteAll = Permission{
		Name:        "user:deleteAll",
		Description: "Allows deleting all user accounts (admin only)",
	}

	UserCreate = Permission{
		Name:        "user:create",
		Description: "Allows creating a new user",
	}
)

var (
	BookCreate = Permission{
		Name:        "book:create",
		Description: "Allows creating a new book",
	}
	BookModify = Permission{
		Name:        "book:modify",
		Description: "Allows modifying book information",
	}
	BookDelete = Permission{
		Name:        "book:delete",
		Description: "Allows deleting a book",
	}
)

var (
	AuthorCreate = Permission{
		Name:        "author:create",
		Description: "Allows creating a new author",
	}
	AuthorModify = Permission{
		Name:        "author:modify",
		Description: "Allows modifying author information",
	}
	AuthorDelete = Permission{
		Name:        "author:delete",
		Description: "Allows deleting an author",
	}
)

var (
	BookEditionCreate = Permission{
		Name:        "bookEdition:create",
		Description: "Allows creating a new book edition",
	}
	BookEditionModify = Permission{
		Name:        "bookEdition:modify",
		Description: "Allows modifying book edition information",
	}
	BookEditionDelete = Permission{
		Name:        "bookEdition:delete",
		Description: "Allows deleting a book edition",
	}
)

var (
	BookTypeCreate = Permission{
		Name:        "bookType:create",
		Description: "Allows creating a new book type",
	}
	BookTypeModify = Permission{
		Name:        "bookType:modify",
		Description: "Allows modifying book type information",
	}
	BookTypeDelete = Permission{
		Name: "bookType:delete",
	}
)

var (
	LocationCreate = Permission{
		Name:        "location:create",
		Description: "Allows creating a new location",
	}
	LocationModify = Permission{
		Name:        "location:modify",
		Description: "Allows modifying location information",
	}
	LocationDelete = Permission{
		Name:        "location:delete",
		Description: "Allows deleting a location",
	}
)
