package middleware

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

// VerifyHasOnePermission checks if the user has at least one of the required permissions from requiredPermissions
// this will return true even if the user only has one of two provided permissions
func VerifyHasOnePermission(requiredPermissions ...permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the user is logged in
		loggedIn, exists := c.Get("loggedIn")
		if !exists || !loggedIn.(bool) {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// get the userId
		userId, exists := c.Get("userId")
		if !exists {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if user exists
		user, err := operations.GetUserById(database.MongoDB, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// get all roles by the user id
		userRoles, err := operations.GetUserRolesByUserId(database.MongoDB, user.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if any role has the required permission
		hasPermission := false
		for _, requiredPermission := range requiredPermissions {
			// get every role and check its permissions
			for _, userRole := range userRoles {
				role, err := operations.GetRoleById(database.MongoDB, userRole.RoleId)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
					c.Abort()
					return
				}

				rolePermissions, err := operations.GetAllRolePermissionsByRoleId(database.MongoDB, role.ID)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
					c.Abort()
					return
				}

				if containsPermission(rolePermissions, requiredPermission.Name) {
					hasPermission = true
					break
				}
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"msg": "Not enough permissions to access this resource", "error": true})
			c.Abort()
			return
		}

	}
}

// VerifyHasAllPermission checks if the user has all the required permissions from requiredPermissions
// this will return true only if the user has all the provided permissions
func VerifyHasAllPermission(requiredPermissions ...permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the user is logged in
		loggedIn, exists := c.Get("loggedIn")
		if !exists || !loggedIn.(bool) {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// get the userId
		userId, exists := c.Get("userId")
		if !exists {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if user exists
		user, err := operations.GetUserById(database.MongoDB, id)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// get all roles by the user id
		userRoles, err := operations.GetUserRolesByUserId(database.MongoDB, user.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if any role has the required permission
		hasPermission := false
		for _, requiredPermission := range requiredPermissions {
			// get every role and check its permissions
			for _, userRole := range userRoles {
				role, err := operations.GetRoleById(database.MongoDB, userRole.RoleId)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
					c.Abort()
					return
				}

				rolePermissions, err := operations.GetAllRolePermissionsByRoleId(database.MongoDB, role.ID)
				if err != nil {
					c.JSON(http.StatusForbidden, gin.H{"msg": "Not enough permissions to access this resource", "error": true})
					c.Abort()
					return
				}

				if containsPermission(rolePermissions, requiredPermission.Name) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{"msg": "Not enough permissions to access this resource", "error": true})
				c.Abort()
				return
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"msg": "Not enough permissions to access this resource", "error": true})
			c.Abort()
			return
		}

	}
}

func containsPermission(slice []models.RolePermission, item string) bool {
	for _, s := range slice {
		if strings.Split(s.Name, "_")[1] == item {
			return true
		}
	}
	return false
}
