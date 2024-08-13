package setup

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/crypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func PerformFirstRunTasks() error {
	// TODO: automatically create all usergroups with the permissions
	// TODO: create the RestrictedUserGroup with the permissions
	// TODO: add first setup frontend page
	//		--> ask for the admin username and password
	//		--> setup the library name
	//		--> ...

	AdminUser, err := createAdminUser()
	if err != nil {
		return err
	}

	NormalUser, err := createNormalUser()
	if err != nil {
		return err
	}

	// create the user permissions
	var userPermissionIds []primitive.ObjectID
	for _, permission := range groups.User.Permissions {
		perm, err := operations.CreatePermission(permission.Name, permission.Description)
		if err != nil {
			return err
		}
		userPermissionIds = append(userPermissionIds, perm.ID)
	}

	// create the admin permissions
	var adminPermissionIds []primitive.ObjectID
	for _, permission := range groups.Administrator.Permissions {
		perm, err := operations.CreatePermission(permission.Name, permission.Description)
		if err != nil {
			return err
		}
		adminPermissionIds = append(adminPermissionIds, perm.ID)
	}

	// create the roles
	normalRole, err := operations.CreateRole(database.MongoDB, groups.User.Name, groups.User.Description, userPermissionIds)
	if err != nil {
		return err
	}
	adminRole, err := operations.CreateRole(database.MongoDB, groups.Administrator.Name, groups.Administrator.Description, adminPermissionIds)
	if err != nil {
		return err
	}

	AdminUser.Roles = append(AdminUser.Roles, adminRole.ID, normalRole.ID)
	NormalUser.Roles = append(NormalUser.Roles, normalRole.ID)

	err = operations.CreateUser(database.MongoDB, AdminUser)
	if err != nil {
		return err
	}

	err = operations.CreateUser(database.MongoDB, NormalUser)
	if err != nil {
		return err
	}

	return nil
}

func createAdminUser() (models.User, error) {
	var adminUser models.User

	randomPW := crypt.GenerateRandomPassword(15, true, true)
	hashedPW, err := crypt.HashPassword(randomPW)
	if err != nil {
		return adminUser, err
	}

	adminUser.ID = primitive.NewObjectID()
	adminUser.Password = hashedPW
	adminUser.Username = "admin"
	// normalUser.Roles =
	adminUser.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	adminUser.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	log.Println("The admin user has been successfully created. The credentials are:")
	log.Println("Username: " + adminUser.Username)
	log.Println("Password: " + randomPW)
	log.Println("Please change the password after your first login !!") // TODO: enforce this

	return adminUser, nil
}

func createNormalUser() (models.User, error) {
	var normalUser models.User

	randomPW := crypt.GenerateRandomPassword(15, true, true)
	hashedPW, err := crypt.HashPassword(randomPW)
	if err != nil {
		return normalUser, err
	}

	normalUser.ID = primitive.NewObjectID()
	normalUser.Password = hashedPW
	normalUser.Username = "testuser"
	// normalUser.Roles =
	normalUser.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	normalUser.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	log.Println("The normal user has been successfully created. The credentials are:")
	log.Println("Username: " + normalUser.Username)
	log.Println("Password: " + randomPW)
	log.Println("Please change the password after your first login !!") // TODO: enforce this

	return normalUser, nil
}
