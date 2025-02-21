package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewareHttp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware is applied only to routes that require authentication.
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		//username, err := validateJWT(tokenString)
		username := "Zack"
		/*

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
		*/

		// Save the username in the context (if needed by the handler).
		c.Set("username", username)
		c.Next()
	}

}
