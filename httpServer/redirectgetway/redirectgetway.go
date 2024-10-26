package redirectgetway

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/medfriend/shared-commons-go/util/consul"
	"io"
	"net/http"
	"os"
	"strings"
)

func Redirectgetway(c *gin.Context) {

	cacheServiceName := os.Getenv("SERVICE_CACHE")

	fmt.Println(cacheServiceName)

	address, port, err := consul.GetServiceAddressAndPort(cacheServiceName)

	if err != nil {
		fmt.Println("cache no se encuentra en consulRegister")
	}

	if err == nil {
		body, errCache := getServiceResponse(c, address, port, cacheServiceName, "GET")

		if body["data"] != "data no avalible on the cache" {
			c.JSON(200, body)
			c.Abort()
			return
		}

		if errCache != nil || body["data"] == "data no avalible on the cache" {

			pathParts := strings.Split(c.Request.URL.Path, "/")

			serviceName := fmt.Sprintf("medfri-%s", strings.Join(pathParts[2:3], "/"))

			address, port, err := consul.GetServiceAddressAndPort(serviceName)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s no se encuentra en consulRegister", serviceName))
			}

			body, errCache = getServiceResponse(c,
				address,
				port,
				serviceName,
				c.Request.Method)

			if body["data"] != "data no avalible on the service" {
				c.JSON(200, body)
				c.Abort()
				return
			}

			return
		}
	}

}

func getServiceResponse(c *gin.Context, address string, port int, serviceGetway string, method string) (map[string]interface{}, error) {

	pathParts := strings.Split(c.Request.URL.Path, "/")
	service := strings.Join(pathParts[2:], "/")

	targetURL := fmt.Sprintf("http://%s:%d/%s/%s", address, port, serviceGetway, service)

	req, err := http.NewRequest(method, targetURL, c.Request.Body)
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
		return nil, fmt.Errorf("Error realizar la solicitud al servicio")
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
