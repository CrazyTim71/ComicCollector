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

func BookTypeHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			bookTypes, err := operations.GetAll[models.BookType](database.Tables.BookType)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			if bookTypes == nil {
				bookTypes = []models.BookType{}
			}

			c.JSON(http.StatusOK, bookTypes)
		})

	rg.GET("/:id",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
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

			bookType, err := operations.GetOneById[models.BookType](database.Tables.BookType, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Book type not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, bookType)
		})

	rg.POST("",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.BookTypeCreate,
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

			// check if the book type already exists
			_, err = operations.GetOneByFilter[models.BookType](database.Tables.BookType, bson.M{"name": requestBody.Name})
			if err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Book type already exists", "error": true})
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

			var newBookType models.BookType
			newBookType.ID = primitive.NewObjectID()
			newBookType.Name = requestBody.Name
			newBookType.Description = requestBody.Description
			newBookType.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
			newBookType.CreatedBy = currentUser

			_, err = operations.InsertOne(database.Tables.BookType, newBookType)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, newBookType)
		})

	rg.PATCH("/:id",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.BookEditionModify,
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

			// check if the book type already exists
			_, err = operations.GetOneById[models.BookType](database.Tables.BookType, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			// update the book type
			result, err := operations.UpdateOne(database.Tables.BookType, objID, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			if result.ModifiedCount == 0 {
				c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Book type was updated successfully"})
		})

	rg.DELETE("/:id",
		middleware.JWTAuth(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.BookTypeDelete,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the book type exists
			_, err = operations.GetOneById[models.BookType](database.Tables.BookType, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Book type doesn't exist", "error": true})
				return
			}

			// delete the book type
			_, err = operations.DeleteOne(database.Tables.BookType, bson.M{"_id": objID})
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Book type was deleted successfully"})
		})
}
