package httpServer

import (
	"fmt"
	"getway-go/httpServer/middleware"
	"getway-go/httpServer/redirectgetway"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func InitHttpServer(taskQueue chan *http.Request) {
	r := gin.Default()

	whitelist := []string{
		"/medfri-getway/security/auth",
	}

	r.Use(middleware.Corsmiddleware())
	r.Use(middleware.Authmiddleware(whitelist))

	r.Any(
		fmt.Sprintf("%s/*path", os.Getenv("SERVICE_PATH")),
		func(c *gin.Context) {
			//taskQueue <- c.Request
			//time.Sleep(time.Second)
			redirectgetway.Redirectgetway(c)
		})

	err := r.Run(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")))

	if err != nil {
		return
	}
}
