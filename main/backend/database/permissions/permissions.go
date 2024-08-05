package permissions

type Permission struct {
	Name        string
	Description string
}

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
	UserDelete = Permission{
		Name:        "user:delete",
		Description: "Allows deleting a user",
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
