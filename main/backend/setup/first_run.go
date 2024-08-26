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

	err = operations.InsertUser(database.MongoDB, AdminUser)
	if err != nil {
		return err
	}

	err = operations.InsertUser(database.MongoDB, NormalUser)
	if err != nil {
		return err
	}

	err = createNoDataEntities()
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

func createNoDataEntities() error {
	// create funcional entities
	var NoAuthor models.Author
	NoAuthor.ID = primitive.NewObjectID()
	NoAuthor.Name = "No Author"
	NoAuthor.Description = "This author is used for books without an author"
	NoAuthor.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	NoAuthor.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err := operations.InsertAuthor(database.MongoDB, NoAuthor)
	if err != nil {
		return err
	}

	var NoPublisher models.Publisher
	NoPublisher.ID = primitive.NewObjectID()
	NoPublisher.Name = "No Publisher"
	NoPublisher.Description = "This publisher is used for books without a publisher"
	NoPublisher.Country = "Narnia"
	NoPublisher.WebsiteURL = "https://narnia.com"
	NoPublisher.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	NoPublisher.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.InsertPublisher(database.MongoDB, NoPublisher)
	if err != nil {
		return err
	}

	var NoLocation models.Location
	NoLocation.ID = primitive.NewObjectID()
	NoLocation.Name = "No Location"
	NoLocation.Description = "This location is used for books without a location"
	NoLocation.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	NoLocation.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.InsertLocation(database.MongoDB, NoLocation)
	if err != nil {
		return err
	}

	var NoOwner models.Owner
	NoOwner.ID = primitive.NewObjectID()
	NoOwner.Name = "No Owner"
	NoOwner.Description = "This owner is used for books without an owner"
	NoOwner.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	NoOwner.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.InsertOwner(database.MongoDB, NoOwner)

	if err != nil {
		return err
	}

	var NoBookEdition models.BookEdition
	NoBookEdition.ID = primitive.NewObjectID()
	NoBookEdition.Name = "No Edition"
	NoBookEdition.Description = "This edition is used for books without an edition"
	NoBookEdition.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	NoBookEdition.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.InsertBookEdition(database.MongoDB, NoBookEdition)
	if err != nil {
		return err
	}

	var NoBookType models.BookType
	NoBookType.ID = primitive.NewObjectID()
	NoBookType.Name = "No Type"
	NoBookType.Description = "This type is used for books without a type"
	NoBookType.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	NoBookType.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.InsertBookType(database.MongoDB, NoBookType)
	if err != nil {
		return err
	}

	return nil
}
