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

func LocationHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			locations, err := operations.GetAllLocations(database.MongoDB)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			if locations == nil {
				locations = []models.Location{}
			}

			c.JSON(http.StatusOK, locations)
		})

	rg.GET("/:id",
		middleware.CheckJwtToken(),
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

			location, err := operations.GetLocationById(database.MongoDB, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Location not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, location)
		})

	rg.POST("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.LocationCreate,
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

			// check if the location already exists
			_, err = operations.GetLocationByName(database.MongoDB, requestBody.Name)
			if err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Location already exists", "error": true})
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

			var newLocation models.Location
			newLocation.ID = primitive.NewObjectID()
			newLocation.Name = requestBody.Name
			newLocation.Description = requestBody.Description
			newLocation.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
			newLocation.CreatedBy = currentUser

			err = operations.InsertLocation(database.MongoDB, newLocation)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, newLocation)
		})

	rg.PATCH("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.LocationModify,
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

			// check if the location already exists
			_, err = operations.GetLocationById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Location not found", "error": true})
				return
			}

			// update the location
			result, err := operations.UpdateLocation(database.MongoDB, objID, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			if result.ModifiedCount == 0 {
				c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Location was updated successfully"})
		})

	rg.DELETE("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.LocationDelete,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the location exists
			_, err = operations.GetLocationById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Location not found", "error": true})
				return
			}

			// delete the location
			_, err = operations.DeleteLocation(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Location was deleted successfully"})
		})
}
