package v1

import (
	"ComicCollector/main/backend/database"
	"ComicCollector/main/backend/database/operations"
	"ComicCollector/main/backend/database/permissions"
	"ComicCollector/main/backend/database/permissions/groups"
	"ComicCollector/main/backend/middleware"
	"ComicCollector/main/backend/utils"
	"ComicCollector/main/backend/utils/env"
	"ComicCollector/main/backend/utils/webcontext"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// image formats and magic numbers
var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

// https://stackoverflow.com/questions/25959386/how-to-check-if-a-file-is-a-valid-image
func mimeFromIncipit(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}

func ImageHandler(rg *gin.RouterGroup) {
	rg.GET("/cover/:id",
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

			buf, err := operations.GetImageById(database.CoverBucket, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Author not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			contentType := http.DetectContentType(buf.Bytes())
			c.Data(http.StatusOK, contentType, buf.Bytes())
		})

	rg.POST("/cover/:bookid",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.CoverUpload,
		),
		middleware.VerifyHasOnePermission(
			permissions.BookCreate,
			permissions.BookModify,
		),
		func(c *gin.Context) {
			var requestBody struct {
				CoverImage *multipart.FileHeader `form:"cover_image" binding:"required"`
			}

			if err := c.ShouldBind(&requestBody); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body", "error": true})
				return
			}

			id := c.Param("bookid")
			bookId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
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
			if !strings.HasSuffix(header.Filename, ".png") && !strings.HasSuffix(header.Filename, ".jpg") && !strings.HasSuffix(header.Filename, ".jpeg") {
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

			// Read the first few bytes of the file to determine its type
			incipit := make([]byte, 8) // Read enough bytes to cover all magic numbers
			_, err = imageFile.Read(incipit)
			if err != nil {
				log.Println("Error: Unable to read file header")
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Invalid file", "error": true})
				return
			}
			// Reset the read pointer again to the start of the file
			_, _ = imageFile.Seek(0, io.SeekStart)

			// Check the file type
			fileType := mimeFromIncipit(incipit)
			if fileType != "image/jpeg" && fileType != "image/png" {
				log.Println("Error: Invalid file type")
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid file type", "error": true})
				return
			}

			// TODO: check if the filename already exists in the db
			uploadStream, err := database.CoverBucket.OpenUploadStream(header.Filename)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to open upload stream", "error": true})
				return
			}
			defer func(uploadStream *gridfs.UploadStream) {
				err := uploadStream.Close()
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to close upload stream", "error": true})
					return
				}
			}(uploadStream)

			_, err = io.Copy(uploadStream, imageFile)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to upload image to GridFS", "error": true})
				return
			}

			fileID := uploadStream.FileID

			// check if the book exists
			_, err = operations.GetBookById(database.MongoDB, bookId)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Book not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			// update the book with the new cover image
			updateData := bson.M{
				"cover_image": fileID,
				"updated_at":  utils.ConvertToDateTime(time.DateTime, time.Now()),
			}

			_, err = operations.UpdateBook(database.MongoDB, bookId, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Cover image uploaded successfully", "error": false})
		})

	// TODO: implement the PATCH method
	//			- delete the old image
	//			- upload the new image
	//			- update the book with the new image
	rg.PATCH("/cover/:bookid")

	rg.DELETE("/cover/:id",
		middleware.CheckJwtToken(),
		middleware.DenyUserGroup(groups.RestrictedUser),
		middleware.VerifyHasAllPermission(
			permissions.BasicApiAccess,
			permissions.CoverDelete,
			permissions.BookModify,
		),
		func(c *gin.Context) {
			id := c.Param("id")
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid ID", "error": true})
				return
			}

			// check if the image exists
			_, err = operations.GetImageById(database.CoverBucket, objID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusNotFound, gin.H{"msg": "Image not found", "error": true})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			// delete the image
			err = operations.DeleteImage(database.CoverBucket, objID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			// remove the image from the book
			book, err := operations.GetBookByCoverImage(database.MongoDB, objID)
			if err != nil {
				// if the book is not found, the image was not associated with any book
				if !errors.Is(err, mongo.ErrNoDocuments) {
					c.JSON(http.StatusOK, gin.H{"msg": "Image deleted successfully", "error": false})
					return
				}
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			currentUser, err := webcontext.GetUserId(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal error", "error": true})
				return
			}

			updateData := bson.M{
				"cover_image": primitive.NilObjectID,
				"updated_at":  utils.ConvertToDateTime(time.DateTime, time.Now()),
				"updated_by":  currentUser,
			}
			_, err = operations.UpdateBook(database.MongoDB, book.ID, updateData)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error", "error": true})
				return
			}

			c.JSON(http.StatusOK, gin.H{"msg": "Image deleted successfully", "error": false})
		})
}
