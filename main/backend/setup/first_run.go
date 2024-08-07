package setup

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/utils/crypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func PerformFirstRunTasks() error {
	// TODO: automatically create all usergroups with the permissions

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
		_, err := createPermission(permission.Name, permission.Description)
		if err != nil {
			return err
		}
	}

	for _, permission := range adminPermissions {
		_, err := createPermission(permission.Name, permission.Description)
		if err != nil {
			return err
		}
	}

	// create the roles
	normalRole, err := createRole(groups.User.Name, groups.User.Description)
	if err != nil {
		return err
	}
	adminRole, err := createRole(groups.Administrator.Name, groups.Administrator.Description)
	if err != nil {
		return err
	}

	// create role permissions
	for _, permission := range normalPermissions {
		perm, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
		if err != nil {
			return err
		}

		_, err = createRolePermission(normalRole, perm)
		if err != nil {
			return err
		}
	}

	for _, permission := range adminPermissions {
		perm, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
		if err != nil {
			return err
		}

		_, err = createRolePermission(adminRole, perm)
		if err != nil {
			return err
		}
	}

	// assign the user roles
	_, err = createUserRole(NormalUser, normalRole)
	if err != nil {
		return err
	}

	// admins have both roles
	_, err = createUserRole(AdminUser, normalRole)
	if err != nil {
		return err
	}

	_, err = createUserRole(AdminUser, adminRole)
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
	//admin_NormalUserRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	//admin_NormalUserRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	//
	//var testuser_NormalUserRole models.UserRole
	//testuser_NormalUserRole.ID = primitive.NewObjectID()
	//testuser_NormalUserRole.UserId = NormalUser.ID
	//testuser_NormalUserRole.RoleId = UserRole.ID
	//testuser_NormalUserRole.Name = NormalUser.Username + "_" + UserRole.Name
	//testuser_NormalUserRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	//testuser_NormalUserRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	//
	//// assigns the administrator role
	//var AdminUserRole models.UserRole
	//AdminUserRole.ID = primitive.NewObjectID()
	//AdminUserRole.UserId = AdminUser.ID
	//AdminUserRole.RoleId = AdminRole.ID
	//AdminUserRole.Name = AdminUser.Username + "_" + AdminRole.Name
	//AdminUserRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	//AdminUserRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
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
	normalUser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	normalUser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	err = operations.SaveUser(database.MongoDB, normalUser)
	if err != nil {
		return normalUser, err
	}

	log.Println("The normal user has been successfully created. The credentials are:")
	log.Println("Username: " + normalUser.Username)
	log.Println("Password: " + randomPW)
	log.Println("Please change the password after your first login !!") // TODO: enforce this

	return normalUser, nil
}

func createRole(name string, description string) (models.Role, error) {
	var role models.Role

	role.ID = primitive.NewObjectID()
	role.Name = name
	role.Description = description
	role.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	role.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	err := operations.SaveRole(database.MongoDB, role)
	if err != nil {
		return role, err
	}

	return role, nil
}

func createPermission(name string, description string) (models.Permission, error) {
	var permission models.Permission

	permission.ID = primitive.NewObjectID()
	permission.Name = name
	permission.Description = description
	permission.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	permission.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// check if permission already exists
	_, err := operations.GetPermissionByName(database.MongoDB, permission.Name)
	if err == nil {
		return permission, nil
	}

	err = operations.SavePermission(database.MongoDB, permission)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func createRolePermission(role models.Role, permission models.Permission) (models.RolePermission, error) {
	var rolePermission models.RolePermission

	rolePermission.ID = primitive.NewObjectID()
	rolePermission.RoleId = role.ID
	rolePermission.PermissionId = permission.ID
	rolePermission.Name = role.Name + "_" + permission.Name
	rolePermission.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	rolePermission.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// check if rolePermission already exists
	_, err := operations.GetRolePermissionByName(database.MongoDB, rolePermission.Name)
	if err == nil {
		return rolePermission, nil
	}

	err = operations.SaveRolePermission(database.MongoDB, rolePermission)
	if err != nil {
		return rolePermission, err
	}

	return rolePermission, nil
}

func createUserRole(user models.User, role models.Role) (models.UserRole, error) {
	var userRole models.UserRole

	userRole.ID = primitive.NewObjectID()
	userRole.UserId = user.ID
	userRole.RoleId = role.ID
	userRole.Name = user.Username + "_" + role.Name
	userRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	userRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// check if userRole already exists
	_, err := operations.GetUserRoleByName(database.MongoDB, userRole.Name)
	if err == nil {
		return userRole, nil
	}

	err = operations.SaveUserRole(database.MongoDB, userRole)
	if err != nil {
		return userRole, err
	}

	return userRole, nil
}

//func createAdminRole() (models.Role, error) {
//	var adminRole models.Role
//
//	adminRole.ID = primitive.NewObjectID()
//	adminRole.Name = "Administrator"
//	adminRole.Description = "Gives a user administrator access"
//	adminRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
//	adminRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
//
//	var viewAllUsersPermission models.Permission
//	viewAllUsersPermission.ID = primitive.NewObjectID()
//	viewAllUsersPermission.Name = permissions.UserViewAll.Name
//	viewAllUsersPermission.Description = permissions.UserViewAll.Description
//	viewAllUsersPermission.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
//	viewAllUsersPermission.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
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
//	userRole.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
//	userRole.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
//
//	var PermissionUserViewSelf models.Permission
//	PermissionUserViewSelf.ID = primitive.NewObjectID()
//	PermissionUserViewSelf.Name = permissions.UserViewSelf.Name
//	PermissionUserViewSelf.Description = permissions.UserViewSelf.Description
//	PermissionUserViewSelf.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
//	PermissionUserViewSelf.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
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
