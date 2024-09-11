package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/webcontext"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func PermissionsHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			allPermissions, err := operations.GetAllPermissions(database.MongoDB)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			if allPermissions == nil {
				allPermissions = []models.Permission{}
			}

			c.JSON(http.StatusOK, allPermissions)
		})

	rg.GET("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			permission, err := operations.GetPermissionById(database.MongoDB, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Permission not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, permission)
		})

	rg.POST("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.PermissionCreate,
		),
		func(c *gin.Context) {
			var requestBody struct {
				Name        string `json:"name" binding:"required"`
				Description string `json:"description" binding:"required"`
			}

			err := c.ShouldBindJSON(&requestBody)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
				return
			}

			// validate the user input
			err = utils.ValidateRequestBody(requestBody)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data. " + err.Error(), "error": true})
				return
			}

			// check if the permission already exists
			_, err = operations.GetPermissionByName(database.MongoDB, requestBody.Name)
			if err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission already exists", "error": true})
				return
			} else if !errors.Is(err, mongo.ErrNoDocuments) {
				// handle all other database errors, but ignore the NoDocuments error
				// that's because this error is expected when the author doesn't exist
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal database error", "error": true})
				return
			}

			currentUser, err := webcontext.GetUserId(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal error", "error": true})
				return
			}

			var newPermission models.Permission
			newPermission.ID = primitive.NewObjectID()
			newPermission.Name = requestBody.Name
			newPermission.Description = requestBody.Description
			newPermission.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
			newPermission.CreatedBy = currentUser

			err = operations.InsertPermission(database.MongoDB, newPermission)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusCreated, newPermission)
		})

	rg.PATCH("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.PermissionModify,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			var requestBody struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}

			if err := c.ShouldBindJSON(&requestBody); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
				return
			}

			// validate the user input
			err = utils.ValidateRequestBody(requestBody)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data. " + err.Error(), "error": true})
				return
			}

			// clean the request body from empty fields to receive only the fields that need to be updated
			updateData := utils.CleanEmptyFields(&requestBody)
			if len(updateData) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "No data provided to update", "error": true})
				return
			}

			currentUser, err := webcontext.GetUserId(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal error", "error": true})
				return
			}
			updateData["updated_at"] = utils.ConvertToDateTime(time.DateTime, time.Now())
			updateData["updated_by"] = currentUser

			// check if the permission already exists
			_, err = operations.GetPermissionById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission not found", "error": true})
				return
			}

			// update the permission
			result, err := operations.UpdatePermission(database.MongoDB, objID, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			if result.ModifiedCount == 0 {
				c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Permission was updated successfully"})
		})

	rg.DELETE("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.PermissionDelete,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the permission exists
			_, err = operations.GetPermissionById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission not found", "error": true})
				return
			}

			// delete the permission
			_, err = operations.DeletePermission(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Permission was deleted successfully"})
		})
}
