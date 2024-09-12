package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func LogRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Log the request body
		log.Println("Request Body:", string(bodyBytes))

		// Replace the request body in the context
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Continue to the next handler
		c.Next()
	}
}
