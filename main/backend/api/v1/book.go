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
		middleware.DenyUserGroup(groups.RestrictedUser), // TODO: test this
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
		),
		func(c *gin.Context) {
			// returns all books
			books, err := operations.GetAllBooks(database.MongoDB)
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
		middleware.DenyUserGroup(groups.RestrictedUser), // TODO: test this
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

			book, err := operations.GetBookById(database.MongoDB, objID)
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

			// TODO: handle the CoverImage file upload
			// TODO: check if BookType, BookEdition, Author, Publisher, Location and Owner are valid

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
			err = utils.ValidateRequestBody(requestBody)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid data. " + err.Error(), "error": true})
				return
			}

			// check if the book already exists
			_, err = operations.GetBookByTitle(database.MongoDB, requestBody.Title)
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

			// TODO: check if the BookType, BookEdition, Author, Publisher, Location and Owner are valid with operations.CheckIfExists()
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
			newBook.UpdatedAt = utils.ConvertToDateTime(time.DateTime, time.Now())

			err = operations.InsertBook(database.MongoDB, newBook)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, newBook)
		})

	//rg.PATCH("")
	//
	//rg.DELETE("")
}
