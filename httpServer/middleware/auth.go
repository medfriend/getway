package middleware

import (
	"fmt"
	"getway-go/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authmiddleware(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("me esta llegando ", c.Request.URL.Path)

		for _, path := range whitelist {
			fmt.Println(path)
			if true {

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
