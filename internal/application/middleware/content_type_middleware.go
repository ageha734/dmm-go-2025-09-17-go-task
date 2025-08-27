package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("ContentTypeMiddleware: Method=%s, Path=%s, ContentType='%s'",
			c.Request.Method, c.Request.URL.Path, c.GetHeader("Content-Type"))

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")

			if contentType == "" {
				log.Printf("ContentTypeMiddleware: Rejecting request with empty Content-Type")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header is required"})
				c.Abort()
				return
			}

			if !strings.HasPrefix(contentType, "application/json") {
				log.Printf("ContentTypeMiddleware: Rejecting request with invalid Content-Type: %s", contentType)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
