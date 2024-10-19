package httpServer

import (
	"fmt"
	"getway-go/httpServer/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

func InitHttpServer() {
	r := gin.Default()

	whitelist := []string{
		"/medfri-getway/test",
	}

	r.Use(middleware.Corsmiddleware())
	r.Use(middleware.Authmiddleware(whitelist))

	r.Any(
		fmt.Sprintf("%s/*path", os.Getenv("SERVICE_PATH")),
		func(c *gin.Context) {

			c.JSON(200, gin.H{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})
		})

	err := r.Run(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))

	if err != nil {
		return
	}
}
