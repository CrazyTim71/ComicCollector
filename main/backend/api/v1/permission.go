package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func PermissionsHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			allPermissions, err := operations.GetAll[models.Permission](database.Tables.Permission)
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
		middleware.JWTAuth(),
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

			permission, err := operations.GetOneById[models.Permission](database.Tables.Permission, objID)
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
		middleware.JWTAuth(),
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
			err = utils.ValidateRequestBody(requestBody, true)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data. " + err.Error(), "error": true})
				return
			}

			// check if the permission already exists
			_, err = operations.GetOneByFilter[models.Permission](database.Tables.Permission, bson.M{"name": requestBody.Name})
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

			currentUser, err := utils.GetUserId(c)
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

			_, err = operations.InsertOne(database.Tables.Permission, newPermission)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusCreated, newPermission)
		})

	rg.PATCH("/:id",
		middleware.JWTAuth(),
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
			err = utils.ValidateRequestBody(requestBody, true)
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

			currentUser, err := utils.GetUserId(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal error", "error": true})
				return
			}
			updateData["updated_at"] = utils.ConvertToDateTime(time.DateTime, time.Now())
			updateData["updated_by"] = currentUser

			// check if the permission already exists
			_, err = operations.GetOneById[models.Permission](database.Tables.Permission, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission not found", "error": true})
				return
			}

			// update the permission
			result, err := operations.UpdateOne(database.Tables.Permission, objID, updateData)
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
		middleware.JWTAuth(),
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
			_, err = operations.GetOneById[models.Permission](database.Tables.Permission, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission not found", "error": true})
				return
			}

			// delete the permission
			_, err = operations.DeleteOne(database.Tables.Permission, bson.M{"_id": objID})
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Permission was deleted successfully"})
		})
}
