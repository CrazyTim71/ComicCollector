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

	normalPermissions := groups.User.Permissions
	adminPermissions := groups.Administrator.Permissions

	// create the permissions
	for _, permission := range normalPermissions {
		_, err := operations.CreatePermission(permission.Name, permission.Description)
		if err != nil {
			return err
		}
	}

	for _, permission := range adminPermissions {
		_, err := operations.CreatePermission(permission.Name, permission.Description)
		if err != nil {
			return err
		}
	}

	// create the roles
	normalRole, err := operations.CreateRole(groups.User.Name, groups.User.Description)
	if err != nil {
		return err
	}
	adminRole, err := operations.CreateRole(groups.Administrator.Name, groups.Administrator.Description)
	if err != nil {
		return err
	}

	// create role permissions
	for _, permission := range normalPermissions {
		perm, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
		if err != nil {
			return err
		}

		_, err = operations.CreateRolePermission(normalRole, perm)
		if err != nil {
			return err
		}
	}

	for _, permission := range adminPermissions {
		perm, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
		if err != nil {
			return err
		}

		_, err = operations.CreateRolePermission(adminRole, perm)
		if err != nil {
			return err
		}
	}

	// assign the user roles
	_, err = operations.CreateUserRole(NormalUser, normalRole)
	if err != nil {
		return err
	}

	// admins have both roles
	_, err = operations.CreateUserRole(AdminUser, normalRole)
	if err != nil {
		return err
	}

	_, err = operations.CreateUserRole(AdminUser, adminRole)
	if err != nil {
		return err
	}

	//AdminRole, err := createAdminRole()
	//if err != nil {
	//	return err
	//}
	//
	//UserRole, err := createUserRole()
	//if err != nil {
	//	return err
	//}
	//
	//// assign the normal user role
	//var admin_NormalUserRole models.UserRole
	//admin_NormalUserRole.ID = primitive.NewObjectID()
	//admin_NormalUserRole.UserId = AdminUser.ID
	//admin_NormalUserRole.RoleId = UserRole.ID
	//admin_NormalUserRole.Name = AdminUser.Username + "_" + UserRole.Name
	//admin_NormalUserRole.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	//admin_NormalUserRole.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	//
	//var testuser_NormalUserRole models.UserRole
	//testuser_NormalUserRole.ID = primitive.NewObjectID()
	//testuser_NormalUserRole.UserId = NormalUser.ID
	//testuser_NormalUserRole.RoleId = UserRole.ID
	//testuser_NormalUserRole.Name = NormalUser.Username + "_" + UserRole.Name
	//testuser_NormalUserRole.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	//testuser_NormalUserRole.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	//
	//// assigns the administrator role
	//var AdminUserRole models.UserRole
	//AdminUserRole.ID = primitive.NewObjectID()
	//AdminUserRole.UserId = AdminUser.ID
	//AdminUserRole.RoleId = AdminRole.ID
	//AdminUserRole.Name = AdminUser.Username + "_" + AdminRole.Name
	//AdminUserRole.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	//AdminUserRole.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	//
	//err = operations.SaveUserRole(database.MongoDB, admin_NormalUserRole)
	//if err != nil {
	//	return err
	//}
	//
	//err = operations.SaveUserRole(database.MongoDB, testuser_NormalUserRole)
	//if err != nil {
	//	return err
	//}
	//
	//err = operations.SaveUserRole(database.MongoDB, AdminUserRole)
	//if err != nil {
	//	return err
	//}

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
	adminUser.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	adminUser.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.CreateUser(database.MongoDB, adminUser)
	if err != nil {
		return adminUser, err
	}

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
	normalUser.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
	normalUser.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

	err = operations.CreateUser(database.MongoDB, normalUser)
	if err != nil {
		return normalUser, err
	}

	log.Println("The normal user has been successfully created. The credentials are:")
	log.Println("Username: " + normalUser.Username)
	log.Println("Password: " + randomPW)
	log.Println("Please change the password after your first login !!") // TODO: enforce this

	return normalUser, nil
}

//func createAdminRole() (models.Role, error) {
//	var adminRole models.Role
//
//	adminRole.ID = primitive.NewObjectID()
//	adminRole.Name = "Administrator"
//	adminRole.Description = "Gives a user administrator access"
//	adminRole.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//	adminRole.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//
//	var viewAllUsersPermission models.Permission
//	viewAllUsersPermission.ID = primitive.NewObjectID()
//	viewAllUsersPermission.Name = permissions.UserViewAll.Name
//	viewAllUsersPermission.Description = permissions.UserViewAll.Description
//	viewAllUsersPermission.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//	viewAllUsersPermission.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//
//	// allows an administrator to view all users
//	var AdminRolePermission models.RolePermission
//	AdminRolePermission.ID = primitive.NewObjectID()
//	AdminRolePermission.Name = adminRole.Name + "_" + viewAllUsersPermission.Name
//	AdminRolePermission.RoleId = adminRole.ID
//	AdminRolePermission.PermissionId = viewAllUsersPermission.ID
//
//	err := operations.SaveRole(database.MongoDB, adminRole)
//	if err != nil {
//		return adminRole, err
//	}
//
//	err = operations.SavePermission(database.MongoDB, viewAllUsersPermission)
//	if err != nil {
//		return adminRole, err
//	}
//
//	err = operations.SaveRolePermission(database.MongoDB, AdminRolePermission)
//	if err != nil {
//		return adminRole, err
//	}
//
//	return adminRole, nil
//}
//
//func createUserRole() (models.Role, error) {
//	var userRole models.Role
//
//	userRole.ID = primitive.NewObjectID()
//	userRole.Name = "User"
//	userRole.Description = "All users are part of this standard group"
//	userRole.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//	userRole.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//
//	var PermissionUserViewSelf models.Permission
//	PermissionUserViewSelf.ID = primitive.NewObjectID()
//	PermissionUserViewSelf.Name = permissions.UserViewSelf.Name
//	PermissionUserViewSelf.Description = permissions.UserViewSelf.Description
//	PermissionUserViewSelf.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//	PermissionUserViewSelf.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
//
//	// assign the permission to the group
//	var RolePermissionUserSelfView models.RolePermission
//	RolePermissionUserSelfView.ID = primitive.NewObjectID()
//	RolePermissionUserSelfView.Name = userRole.Name + "_" + PermissionUserViewSelf.Name
//	RolePermissionUserSelfView.RoleId = userRole.ID
//	RolePermissionUserSelfView.PermissionId = PermissionUserViewSelf.ID
//
//	err := operations.SaveRole(database.MongoDB, userRole)
//	if err != nil {
//		return userRole, err
//	}
//
//	err = operations.SavePermission(database.MongoDB, PermissionUserViewSelf)
//	if err != nil {
//		return userRole, err
//	}
//
//	err = operations.SaveRolePermission(database.MongoDB, RolePermissionUserSelfView)
//	if err != nil {
//		return userRole, err
//	}
//
//	return userRole, nil
//}
