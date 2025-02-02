package middleware

import (
	"getway-go/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authmiddleware(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		for _, path := range whitelist {
			if path == c.Request.URL.Path {
				c.Next()
				return
			}
		}

		authHeader := c.Request.Header.Get("Authorization")

		token, err := jwt.ValidateJWT(authHeader)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido"})
			c.Abort()
			return
		}

		c.Next()
	}
}
