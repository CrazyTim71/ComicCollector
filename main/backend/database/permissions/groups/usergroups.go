package groups

import "ComicCollector/main/backend/database/permissions"

type UserGroup struct {
	Name        string
	Description string
	Permissions []permissions.Permission
}

var (
	Administrator = UserGroup{
		Name:        "Administrator",
		Description: "Full access to all features",
		Permissions: []permissions.Permission{
			// user stuff
			permissions.UserViewAll,
			permissions.UserModifyAll,
			permissions.UserDeleteAll,
			permissions.UserCreate,

			// book stuff
			permissions.BookCreate,
			permissions.BookModify,
			permissions.BookDelete,
		},
	}

	User = UserGroup{
		Name:        "User",
		Description: "Basic access to all features",
		Permissions: []permissions.Permission{
			// user stuff
			permissions.UserViewSelf,
			permissions.UserModifySelf,
			permissions.UserDeleteSelf,
		},
	}

	RestrictedUser = UserGroup{
		Name:        "RestrictedUser",
		Description: "Restricted access to all features until an administrator approves the user",
		Permissions: []permissions.Permission{
			// user stuff
			permissions.UserViewSelf,
		},
	}
)
