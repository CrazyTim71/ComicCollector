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

func CheckUserGroup(userID primitive.ObjectID, group UserGroup) (bool, error) {
	// check if user exists
	user, err := operations.GetUserById(database.MongoDB, userID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
		}
		return false, err
	}

	// get all roles by the user id
	userRoles, err := operations.GetUserRolesByUserId(database.MongoDB, user.ID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
		}
		return false, err
	}

	// check if the user has the admin role
	isUserInGroup := false
	for _, userRole := range userRoles {
		// get single role
		role, err := operations.GetRoleById(database.MongoDB, userRole.RoleId)
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
