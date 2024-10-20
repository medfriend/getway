package redirectgetway

import (
	"fmt"
	"getway-go/consul"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

func Redirectgetway(c *gin.Context) {

	address, port, err := consul.GetServiceAddressAndPort(os.Getenv("SERVICE_CACHE"))

	handlerCache(c, address, port)

	if err != nil {
		fmt.Println("cache no se encuentra en consul")
	}

	c.JSON(200, gin.H{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	})
}

func handlerCache(c *gin.Context, address string, port int) {

	pathParts := strings.Split(c.Request.URL.Path, "/")
	service := strings.Join(pathParts[2:], "/")

	targetURL := fmt.Sprintf("http://%s:%d/medfri-cache/%s", address, port, service)

	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		fmt.Println("Error creando la nueva solicitud", err)
		c.JSON(500, gin.H{"error": "Error creando la solicitud al servicio de cache"})
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la solicitud al servicio de cache", err)
		c.JSON(500, gin.H{"error": "No se pudo comunicar con el servicio de cache"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta del servicio de cache", err)
		c.JSON(500, gin.H{"error": "Error al leer la respuesta del servicio de cache"})
		return
	}

	fmt.Println(string(body))

}
