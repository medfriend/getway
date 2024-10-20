package redirectgetway

import (
	"encoding/json"
	"fmt"
	"getway-go/util/consul"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

func Redirectgetway(c *gin.Context) {

	address, port, err := consul.GetServiceAddressAndPort(os.Getenv("SERVICE_CACHE"))

	if err != nil {
		fmt.Println("cache no se encuentra en consul")
	}
	if err == nil {
		body, errCache := getCache(c, address, port)

		if body["data"] != "data no avalible on the cache" {
			c.JSON(200, body)
			c.Abort()
			return
		}

		if errCache != nil || body["data"] == "data no avalible on the cache" {
			fmt.Println("entonces mi gente")
			return
		}
	}

}

func getCache(c *gin.Context, address string, port int) (map[string]interface{}, error) {

	pathParts := strings.Split(c.Request.URL.Path, "/")
	service := strings.Join(pathParts[2:], "/")

	targetURL := fmt.Sprintf("http://%s:%d/medfri-cache/%s", address, port, service)

	req, err := http.NewRequest("GET", targetURL, c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("Error creando la nueva solicitud")
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error realizar la solicitud al servicio de cache")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error al leer la respuesta del servicio de cache")
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta del servicio de cache")
	}

	return result, nil
}
