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

func RoleHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			roles, err := operations.GetAll[models.Role](database.Tables.Role)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			if roles == nil {
				roles = []models.Role{}
			}

			c.JSON(http.StatusOK, roles)
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

			role, err := operations.GetOneById[models.Role](database.Tables.Role, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Role not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, role)
		})

	rg.POST("",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.RoleCreate,
		),
		func(c *gin.Context) {
			var requestBody struct {
				Name        string               `json:"name" binding:"required"`
				Description string               `json:"description" binding:"required"`
				Permissions []primitive.ObjectID `json:"permissions" binding:"required"`
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

			// check if the role already exists
			_, err = operations.GetOneByFilter[models.Role](database.Tables.Role, bson.M{"name": requestBody.Name})
			if err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Role already exists", "error": true})
				return
			} else if !errors.Is(err, mongo.ErrNoDocuments) {
				// handle all other database errors, but ignore the NoDocuments error
				// that's because this error is expected when the author doesn't exist
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal database error", "error": true})
				return
			}

			// check if the permissions exist
			if !operations.CheckIfAllIdsExist[models.Permission](database.Tables.Permission, requestBody.Permissions) {
				log.Println("Not all provided permission ids exist/are valid")
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Not all provided permission ids exist/are valid", "error": true})
				return
			}

			currentUser, err := utils.GetUserId(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal error", "error": true})
				return
			}

			var newRole models.Role
			newRole.ID = primitive.NewObjectID()
			newRole.Name = requestBody.Name
			newRole.Description = requestBody.Description
			newRole.Permissions = requestBody.Permissions
			newRole.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
			newRole.CreatedBy = currentUser

			_, err = operations.InsertOne(database.Tables.Role, newRole)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, newRole)
		})

	rg.PATCH("/:id",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.RoleModify,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			var requestBody struct {
				Name        string               `json:"name"`
				Description string               `json:"description"`
				Permissions []primitive.ObjectID `json:"permissions"`
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

			// check if the role already exists
			_, err = operations.GetOneById[models.Role](database.Tables.Role, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Role not found", "error": true})
				return
			}

			// update the role
			result, err := operations.UpdateOne(database.Tables.Role, objID, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			if result.ModifiedCount == 0 {
				c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Role was updated successfully"})
		})

	rg.DELETE("/:id",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyUserGroup(groups.Administrator),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.RoleDelete,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the role exists
			_, err = operations.GetOneById[models.Role](database.Tables.Role, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Role not found", "error": true})
				return
			}

			// delete the role
			_, err = operations.DeleteOne(database.Tables.Role, bson.M{"_id": objID})
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Role was deleted successfully"})
		})
}
