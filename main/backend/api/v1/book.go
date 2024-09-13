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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

// BookHandler api/v1/book
func BookHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			// returns all books
			books, err := operations.GetAll[models.Book](database.Tables.Book)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			if books == nil {
				books = []models.Book{}
			}

			c.JSON(http.StatusOK, books)
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

			book, err := operations.GetOneById[models.Book](database.Tables.Book, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Author not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, book)
		})

	rg.POST("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.BookCreate,
		),
		func(c *gin.Context) {
			var requestBody struct {
				Title       string `json:"title" binding:"required"`
				Number      int    `json:"number" binding:"required"`
				ReleaseDate string `json:"release_date" binding:"required"`
				//CoverImage  *multipart.FileHeader `form:"cover_image" binding:"required"`
				Description string               `json:"description"`
				Notes       string               `json:"notes"`
				Authors     []primitive.ObjectID `json:"authors"`
				Publishers  []primitive.ObjectID `json:"publishers"`
				Locations   []primitive.ObjectID `json:"locations"`
				Owners      []primitive.ObjectID `json:"owners" binding:"required"`
				BookType    primitive.ObjectID   `json:"book_type" binding:"required"`
				BookEdition primitive.ObjectID   `json:"book_edition" binding:"required"`
				Printing    string               `json:"printing"`
				ISBN        string               `json:"isbn"`
				Price       string               `json:"price"`
				Count       int                  `json:"count"`
			}

			if err := c.ShouldBindJSON(&requestBody); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
				return
			}

			date, err := time.Parse(time.DateOnly, requestBody.ReleaseDate) // Format must match the date string in the frontend
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
				return
			}
			releaseDate := utils.ConvertToDateTime(time.DateOnly, date)

			// validate the user input
			err = utils.ValidateRequestBody(requestBody, true)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data. " + err.Error(), "error": true})
				return
			}

			// check if the book already exists
			_, err = operations.GetOneByFilter[models.Book](database.Tables.Book, bson.M{"title": requestBody.Title, "number": requestBody.Number})
			if err == nil { // err == nil in case the book already exists
				log.Println(err)
				c.JSON(http.StatusConflict, gin.H{"msg": "This book already exists", "error": true})
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

			var newBook models.Book
			newBook.ID = primitive.NewObjectID()
			newBook.Title = requestBody.Title
			newBook.Number = requestBody.Number
			newBook.ReleaseDate = releaseDate
			newBook.CoverImage = primitive.NilObjectID // the cover image will be uploaded separately
			newBook.Description = requestBody.Description
			newBook.Notes = requestBody.Notes
			newBook.Authors = requestBody.Authors
			newBook.Publishers = requestBody.Publishers
			newBook.Locations = requestBody.Locations
			newBook.Owners = requestBody.Owners
			newBook.BookType = requestBody.BookType
			newBook.BookEdition = requestBody.BookEdition
			newBook.Printing = requestBody.Printing
			newBook.ISBN = requestBody.ISBN
			newBook.Price = requestBody.Price
			newBook.Count = requestBody.Count
			newBook.CreatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())
			newBook.CreatedBy = currentUser

			_, err = operations.InsertOne(database.Tables.Book, newBook)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, newBook)
		})

	rg.PATCH("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.BookModify,
		),
		func(c *gin.Context) {
			id := c.Param("id")

			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			var requestBody struct {
				Title       string               `json:"title"`
				Number      int                  `json:"number"`
				ReleaseDate string               `json:"release_date"`
				Description string               `json:"description"`
				Notes       string               `json:"notes"`
				Authors     []primitive.ObjectID `json:"authors"`
				Publishers  []primitive.ObjectID `json:"publishers"`
				Locations   []primitive.ObjectID `json:"locations"`
				Owners      []primitive.ObjectID `json:"owners"`
				BookType    primitive.ObjectID   `json:"book_type"`
				BookEdition primitive.ObjectID   `json:"book_edition"`
				Printing    string               `json:"printing"`
				ISBN        string               `json:"isbn"`
				Price       string               `json:"price"`
				Count       int                  `json:"count"`
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

			currentUser, err := webcontext.GetUserId(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal error", "error": true})
				return
			}
			updateData["updated_at"] = utils.ConvertToDateTime(time.DateTime, time.Now())
			updateData["updated_by"] = currentUser

			// check if the book exists
			_, err = operations.GetOneById[models.Book](database.Tables.Book, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusNotFound, gin.H{"msg": "Book not found", "error": true})
				return
			}

			result, err := operations.UpdateOne(database.Tables.Book, objID, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			if result.ModifiedCount == 0 {
				c.JSON(http.StatusNotModified, gin.H{"msg": "Nothing was updated", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Book updated successfully", "error": false})
		})

	rg.DELETE("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.BookDelete,
		),
		func(c *gin.Context) {
			id := c.Param("id")
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the book exists
			book, err := operations.GetOneById[models.Book](database.Tables.Book, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Book not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			// delete the cover image
			if book.CoverImage != primitive.NilObjectID {
				err = operations.DeleteImage(database.CoverBucket, book.CoverImage)
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error: unable to delete the old cover file", "error": true})
					return
				}
			}

			// delete the book
			_, err = operations.DeleteOne(database.Tables.Book, bson.M{"_id": objID})
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Book deleted successfully", "error": false})
		})
}
