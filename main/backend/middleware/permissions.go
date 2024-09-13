package middleware

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/utils/webcontext"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
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
		userId, err := webcontext.GetUserId(c)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// check if user exists
		user, err := operations.GetOneById[models.User](database.Tables.User, userId)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				log.Println(err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if any role has the required permission
		hasPermission := false
		for _, requiredPermission := range requiredPermissions {
			// get every role and check its permissions
			for _, roleId := range user.Roles {
				rolePermissions, err := operations.GetAllPermissionsFromRole(database.MongoDB, roleId)
				if err != nil {
					if !errors.Is(err, mongo.ErrNoDocuments) {
						log.Println(err)
					}
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

		userId, err := webcontext.GetUserId(c)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// check if user exists
		user, err := operations.GetOneById[models.User](database.Tables.User, userId)
		if err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				log.Println(err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		// check if any role has the required permission
		hasPermission := false
		for _, requiredPermission := range requiredPermissions {
			// get every role and check its permissions
			for _, roleId := range user.Roles {
				rolePermissions, err := operations.GetAllPermissionsFromRole(database.MongoDB, roleId)
				if err != nil {
					if !errors.Is(err, mongo.ErrNoDocuments) {
						log.Println(err)
					}
					c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
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

func containsPermission(slice []models.Permission, item string) bool {
	for _, s := range slice {
		if s.Name == item {
			return true
		}
	}
	return false
}

func VerifyUserGroup(group groups.UserGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the user is logged in
		loggedIn, exists := c.Get("loggedIn")
		if !exists || !loggedIn.(bool) {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userId, err := webcontext.GetUserId(c)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		isGroup, err := groups.CheckUserGroup(userId, group)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		if isGroup {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "Not enough permissions to view this site", "error": true})
			c.Abort()
			return
		}
	}
}

func DenyUserGroup(group groups.UserGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the user is logged in
		loggedIn, exists := c.Get("loggedIn")
		if !exists || !loggedIn.(bool) {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userId, err := webcontext.GetUserId(c)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		isGroup, err := groups.CheckUserGroup(userId, group)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": true})
			c.Abort()
			return
		}

		if isGroup {
			c.JSON(http.StatusForbidden, gin.H{"message": "Not enough permissions to view this site", "error": true})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
