package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/models"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/env"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// BookHandler api/v1/book
func BookHandler(rg *gin.RouterGroup) {
	rg.GET("",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser), // TODO: test this
		middleware.VerifyHasAllPermission(
			permissions.BookViewAll,
		),
		func(c *gin.Context) {
			// returns all books
			books, err := operations.GetAllBooks(database.MongoDB)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
			c.JSON(http.StatusOK, books)
		})

	rg.GET("/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser), // TODO: test this
		middleware.VerifyHasAllPermission(
			permissions.BookViewAll,
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
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, book)
		})

	rg.POST("",
		middleware.CheckJwtToken(),
		middleware.VerifyHasAllPermission(
			permissions.BookCreate,
		),
		func(c *gin.Context) {

			// create a new book
			var requestBody struct {
				// ID          primitive.ObjectID `form:"id" bson:"_id"`
				Title       string             `form:"title" binding:"required"`
				Number      int                `form:"number" binding:"required"`
				ReleaseDate string             `form:"release_date" binding:"required"`
				Description string             `form:"description" binding:"required"`
				Notes       string             `form:"notes" binding:"required"`
				BookType    primitive.ObjectID `form:"book_type" binding:"required"`
				BookEdition primitive.ObjectID `form:"book_edition" binding:"required"`
				Printing    string             `form:"printing" binding:"required"`
				ISBN        string             `form:"isbn" binding:"required"`
				Price       string             `form:"price" binding:"required"`
				Count       int                `form:"count" binding:"required"`

				Author    primitive.ObjectID `form:"author" binding:"required"`
				Publisher primitive.ObjectID `form:"publisher" binding:"required"`
				Location  primitive.ObjectID `form:"location" binding:"required"`
				Owner     primitive.ObjectID `form:"owner" binding:"required"`

				// CoverImage  string             `form:"cover_image" binding:"required"`
				// CreatedAt   primitive.DateTime `form:"created_at" bson:"created_at"`
				// UpdatedAt   primitive.DateTime `form:"updated_at" bson:"updated_at"`
			}

			// TODO: verify the user input with joi
			// TODO: handle the CoverImage file upload
			// TODO: check if BookType, BookEdition, Author, Publisher, Location and Owner are valid

			if err := c.ShouldBindJSON(&requestBody); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
				return
			}

			// Handle the CoverImage file upload
			imageFile, header, err := c.Request.FormFile("cover_image")
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "The cover image is required", "error": true})
				return
			}

			// Check the file extension
			if !strings.HasSuffix(header.Filename, ".png") && !strings.HasSuffix(header.Filename, ".jpg") {
				log.Println("Invalid file extension")
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid file", "error": true})
				return
			}

			// Check the file size
			size, err := imageFile.Seek(0, io.SeekEnd)
			if err != nil {
				log.Println("Error: Unable to determine file size")
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Invalid file", "error": true})
				return
			}
			// Reset the read pointer to the start of the file
			_, _ = imageFile.Seek(0, io.SeekStart)

			if size > int64(env.MaxImageFileSize) {
				log.Println("Error: File size exceeds the limit")
				c.JSON(http.StatusBadRequest, gin.H{"msg": "File size exceeds limit of \" + fmt.Sprint(env.MaxImageFileSize>>20) + \" MiB", "error": true})
				return
			}

			//// create upload folder if it doesn't exist
			//uploadDir := "./uploads"
			//if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			//	if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			//		log.Println(err)
			//		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Unable to save the file", "error": true})
			//		return
			//	}
			//}
			//
			//// create a unique filename
			//uniqueFilename := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filepath.Ext(header.Filename)
			//filePath := filepath.Join(uploadDir, uniqueFilename)
			//
			//// save file to disk into the temp folder
			//if err := c.SaveUploadedFile(header, filePath); err != nil {
			//	log.Println(err)
			//	c.JSON(http.StatusInternalServerError, gin.H{"msg": "Unable to save file", "error": true})
			//	return
			//}

			imageData, err := io.ReadAll(imageFile)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to read cover image", "error": true})
				return
			}

			newBook := models.Book{
				ID:          primitive.NewObjectID(),
				Title:       requestBody.Title,
				Number:      requestBody.Number,
				ReleaseDate: requestBody.ReleaseDate,
				CoverImage:  imageData,
				Description: requestBody.Description,
				Notes:       requestBody.Notes,
				BookType:    requestBody.BookType,
				BookEdition: requestBody.BookEdition,
				Printing:    requestBody.Printing,
				ISBN:        requestBody.ISBN,
				Price:       requestBody.Price,
				Count:       requestBody.Count,
				CreatedAt:   utils.ConvertToDateTime(time.DateTime, time.Now()),
				UpdatedAt:   utils.ConvertToDateTime(time.DateTime, time.Now()),
			}

			err = operations.CreateBook(database.MongoDB, newBook)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}
		})

	//rg.PATCH("")
	//
	//rg.DELETE("")
}
