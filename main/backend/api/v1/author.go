package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/JoiHelper"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func AuthorHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		func(c *gin.Context) {
			authors, err := operations.GetAllAuthors(database.MongoDB)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			if authors == nil {
				authors = []models.Author{}
			}

			c.JSON(http.StatusOK, authors)
		})

	rg.GET("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			author, err := operations.GetAuthorById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, author)
		})

	rg.POST("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.AuthorCreate,
		),
		func(c *gin.Context) {
			var requestBody struct {
				Name        string `json:"name" binding:"required"`
				Description string `json:"description" binding:"required"`
			}

			if err := c.ShouldBindJSON(&requestBody); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
				return
			}

			//validate the user input
			err := JoiHelper.UserInput.Validate(requestBody.Name)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid name. Please remove all invalid characters and try again.", "error": true})
				return
			}

			// check if the author already exists
			_, err = operations.GetAuthorByName(database.MongoDB, requestBody.Name)
			if err == nil { // err == nil in case the author already exists
				log.Println(err)
				c.JSON(http.StatusConflict, gin.H{"msg": "This author already exists", "error": true})
				return
			} else if !errors.Is(err, mongo.ErrNoDocuments) {
				// handle all other database errors, but ignore the NoDocuments error
				// that's because this error is expected when the author doesn't exist
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal database error", "error": true})
				return
			}

			var newAuthor models.Author
			newAuthor.ID = primitive.NewObjectID()
			newAuthor.Name = requestBody.Name
			newAuthor.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
			newAuthor.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

			err = operations.CreateAuthor(database.MongoDB, newAuthor)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Author was created successfully"})
		})

	rg.PATCH("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.AuthorModify,
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

			//validate the user input
			// TODO: validate all fields
			for _, v := range requestBody {

			}
			err = JoiHelper.UserInput.Validate(requestBody.Name)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid name. Please remove all invalid characters and try again.", "error": true})
				return
			}

			// check if the author exists
			existingAuthor, err := operations.GetAuthorById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusNotFound, gin.H{"msg": "Author not found", "error": true})
				return
			}

			// update the author
			// TODO: check and update all fields
			existingAuthor.Name = requestBody.Name

			err = operations.UpdateAuthor(database.MongoDB, existingAuthor)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Author was updated successfully"})
		})

	rg.DELETE("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.AuthorDelete,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the author exists
			_, err = operations.GetAuthorById(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusNotFound, gin.H{"msg": "Author not found", "error": true})
				return
			}

			err = operations.DeleteAuthor(database.MongoDB, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Author was deleted successfully"})
		})
}
