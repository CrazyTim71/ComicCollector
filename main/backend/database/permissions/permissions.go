package permissions

type Permission struct {
	name        string
	description string
}

func (p Permission) Name() string {
	return p.name
}

func (p Permission) Description() string {
	return p.description
}

var (
	BasicApiAccess = Permission{
		name:        "basic:ApiAccess",
		description: "Gives an user the ability to request data from the api endpoints",
	}
)

var (
	UserViewSelf = Permission{
		name:        "user:viewSelf",
		description: "Allows a user to view their own information",
	}
	UserViewAll = Permission{
		name:        "user:viewAll",
		description: "Allows viewing all user information (admin only)",
	}

	UserModifySelf = Permission{
		name:        "user:modifySelf",
		description: "Allows a user to modify their own user information",
	}
	UserModifyAll = Permission{
		name:        "user:modifyAll",
		description: "Allows modifying all user information (admin only)",
	}

	UserDeleteSelf = Permission{
		name:        "user:deleteSelf",
		description: "Allows a user to delete their own account",
	}
	UserDeleteAll = Permission{
		name:        "user:deleteAll",
		description: "Allows deleting all user accounts (admin only)",
	}

	UserCreate = Permission{
		name:        "user:create",
		description: "Allows creating a new user",
	}
)

var (
	BookCreate = Permission{
		name:        "book:create",
		description: "Allows creating a new book",
	}
	BookModify = Permission{
		name:        "book:modify",
		description: "Allows modifying book information",
	}
	BookDelete = Permission{
		name:        "book:delete",
		description: "Allows deleting a book",
	}
)

var (
	AuthorCreate = Permission{
		name:        "author:create",
		description: "Allows creating a new author",
	}
	AuthorModify = Permission{
		name:        "author:modify",
		description: "Allows modifying author information",
	}
	AuthorDelete = Permission{
		name:        "author:delete",
		description: "Allows deleting an author",
	}
)

var (
	BookEditionCreate = Permission{
		name:        "bookEdition:create",
		description: "Allows creating a new book edition",
	}
	BookEditionModify = Permission{
		name:        "bookEdition:modify",
		description: "Allows modifying book edition information",
	}
	BookEditionDelete = Permission{
		name:        "bookEdition:delete",
		description: "Allows deleting a book edition",
	}
)

var (
	BookTypeCreate = Permission{
		name:        "bookType:create",
		description: "Allows creating a new book type",
	}
	BookTypeModify = Permission{
		name:        "bookType:modify",
		description: "Allows modifying book type information",
	}
	BookTypeDelete = Permission{
		name: "bookType:delete",
	}
)

var (
	LocationCreate = Permission{
		name:        "location:create",
		description: "Allows creating a new location",
	}
	LocationModify = Permission{
		name:        "location:modify",
		description: "Allows modifying location information",
	}
	LocationDelete = Permission{
		name:        "location:delete",
		description: "Allows deleting a location",
	}
)

var (
	OwnerCreate = Permission{
		name:        "owner:create",
		description: "Allows creating a new owner",
	}
	OwnerModify = Permission{
		name:        "owner:modify",
		description: "Allows modifying owner information",
	}
	OwnerDelete = Permission{
		name:        "owner:delete",
		description: "Allows deleting an owner",
	}
)

var (
	PermissionCreate = Permission{
		name:        "permission:create",
		description: "Allows creating a new permission",
	}
	PermissionModify = Permission{
		name:        "permission:modify",
		description: "Allows modifying permission information",
	}
	PermissionDelete = Permission{
		name:        "permission:delete",
		description: "Allows deleting a permission",
	}
)

var (
	PublisherCreate = Permission{
		name:        "publisher:create",
		description: "Allows creating a new publisher",
	}
	PublisherModify = Permission{
		name:        "publisher:modify",
		description: "Allows modifying publisher information",
	}
	PublisherDelete = Permission{
		name:        "publisher:delete",
		description: "Allows deleting a publisher",
	}
)

var (
	RoleCreate = Permission{
		name:        "role:create",
		description: "Allows creating a new role",
	}
	RoleModify = Permission{
		name:        "role:modify",
		description: "Allows modifying role information",
	}
	RoleDelete = Permission{
		name:        "role:delete",
		description: "Allows deleting a role",
	}
)

var (
	CoverUpload = Permission{
		name:        "cover:upload",
		description: "Allows uploading a cover image",
	}
	CoverDelete = Permission{
		name:        "cover:delete",
		description: "Allows deleting a cover image",
	}
)
