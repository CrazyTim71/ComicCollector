package setup

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/utils/crypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func PerformFirstRunTasks() error {
	// TODO: create basic permissions, like admin and user

	AdminUser, err := createAdminUser()
	if err != nil {
		return err
	}

	err = createAdminRoleAndPermissions(AdminUser)
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
	adminUser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	adminUser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	err = operations.SaveUser(database.MongoDB, adminUser)
	if err != nil {
		return adminUser, err
	}

	log.Println("The admin user has been successfully created. The credentials are:")
	log.Println("Username: " + adminUser.Username)
	log.Println("Password: " + randomPW)
	log.Println("Please change the password after your first login !!") // TODO: enforce this

	return adminUser, nil
}

func createAdminRoleAndPermissions(adminUser models.User) error {
	var adminRole models.Role

	adminRole.ID = primitive.NewObjectID()
	adminRole.Name = "Administrator"
	adminRole.Description = "Gives a user administrator access"
	adminRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	adminRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	var createBookPermission models.Permission
	createBookPermission.ID = primitive.NewObjectID()
	createBookPermission.Name = "createBook"
	createBookPermission.Description = "Create a new book"
	createBookPermission.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	createBookPermission.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// allows an administrator to create books
	var AdminRolePermission models.RolePermission
	AdminRolePermission.ID = primitive.NewObjectID()
	AdminRolePermission.Name = adminRole.Name + "_" + createBookPermission.Name
	AdminRolePermission.RoleId = adminRole.ID
	AdminRolePermission.PermissionId = createBookPermission.ID

	// assigns the createBookPermission to the administrator role
	var AdminUserRole models.UserRole
	AdminUserRole.ID = primitive.NewObjectID()
	AdminUserRole.UserId = adminUser.ID
	AdminUserRole.RoleId = adminRole.ID
	AdminUserRole.Name = adminUser.Username + "_" + adminRole.Name
	AdminUserRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	AdminUserRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := operations.SaveRole(database.MongoDB, adminRole)
	if err != nil {
		return err
	}

	err = operations.SavePermission(database.MongoDB, createBookPermission)
	if err != nil {
		return err
	}

	err = operations.SaveRolePermission(database.MongoDB, AdminRolePermission)
	if err != nil {
		return err
	}

	err = operations.SaveUserRole(database.MongoDB, AdminUserRole)
	if err != nil {
		return err
	}

	return nil
}
