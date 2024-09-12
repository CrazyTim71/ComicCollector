package groups

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

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
			//global stuff
			permissions.BasicApiAccess,

			// user stuff
			permissions.UserViewAll,
			permissions.UserModifyAll,
			permissions.UserDeleteAll,
			permissions.UserCreate,

			// book stuff
			permissions.BookCreate,
			permissions.BookModify,
			permissions.BookDelete,

			// author stuff
			permissions.AuthorCreate,
			permissions.AuthorModify,
			permissions.AuthorDelete,

			// book edition stuff
			permissions.BookEditionCreate,
			permissions.BookEditionModify,
			permissions.BookEditionDelete,

			// book type stuff
			permissions.BookTypeCreate,
			permissions.BookTypeModify,
			permissions.BookTypeDelete,

			// location stuff
			permissions.LocationCreate,
			permissions.LocationModify,
			permissions.LocationDelete,

			// owner stuff
			permissions.OwnerCreate,
			permissions.OwnerModify,
			permissions.OwnerDelete,

			// permission stuff
			permissions.PermissionCreate,
			permissions.PermissionModify,
			permissions.PermissionDelete,

			// publisher stuff
			permissions.PublisherCreate,
			permissions.PublisherModify,
			permissions.PublisherDelete,

			// role stuff
			permissions.RoleCreate,
			permissions.RoleModify,
			permissions.RoleDelete,

			// cover stuff
			permissions.CoverUpload,
			permissions.CoverDelete,
		},
	}

	User = UserGroup{
		Name:        "User",
		Description: "Basic access to all features",
		Permissions: []permissions.Permission{
			//global stuff
			permissions.BasicApiAccess,

			// user stuff
			permissions.UserViewSelf,
			permissions.UserModifySelf,
			permissions.UserDeleteSelf,

			// book stuff
			// no permissions yet

			// author stuff
			permissions.AuthorCreate,

			// book edition stuff
			permissions.BookEditionCreate,

			// book type stuff
			permissions.BookTypeCreate,

			// location stuff
			permissions.LocationCreate,

			// owner stuff
			permissions.OwnerCreate,

			// permission stuff
			// no permissions yet

			// publisher stuff
			permissions.PublisherCreate,

			// role stuff
			// no permissions yet

			// cover stuff
			permissions.CoverUpload,
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

func CheckUserGroup(userID primitive.ObjectID, group UserGroup) (bool, error) {
	// check if user exists
	user, err := operations.GetUserById(database.MongoDB, userID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
		}
		return false, err
	}

	// check if the user has the desired role
	isUserInGroup := false
	for _, roleId := range user.Roles {
		// get single role
		role, err := operations.GetRoleById(database.MongoDB, roleId)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				log.Println(err)
			}
			return false, err
		}

		if role.Name == group.Name {
			isUserInGroup = true
			break
		}
	}

	return isUserInGroup, nil
}
