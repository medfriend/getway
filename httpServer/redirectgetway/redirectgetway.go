package redirectgetway

import (
	"fmt"
	"getway-go/httpServer/service"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/medfriend/shared-commons-go/util/consul"
	"os"
	"strings"
)

func Redirectgetway(c *gin.Context, consulClient *api.Client) {

	ignoreCache := c.Request.Header.Get("ignore-cache")

	cacheServiceName := os.Getenv("SERVICE_CACHE")

	address, port, err := consul.GetServiceAddressAndPort(consulClient, cacheServiceName)

	if ignoreCache == "Y" {
		registerOnService(c, address, port, cacheServiceName, consulClient)
		return
	}

	if err == nil {

		body, errCache, cacheStatusCode := service.GetServiceResponse(c, address, port, cacheServiceName, "GET", true)

		if body["data"] != "data no avalible on the cache" {
			c.JSON(*cacheStatusCode, body)
			c.Abort()
			return
		}

		if errCache != nil || body["data"] == "data no avalible on the cache" {
			registerOnService(c, address, port, cacheServiceName, consulClient)
		}
	}

}

func registerOnService(c *gin.Context, address string, port int, cacheServiceName string, consulClient *api.Client) {

	ignoreCache := c.Request.Header.Get("ignore-cache")
	pathParts := strings.Split(c.Request.URL.Path, "/")
	serviceName := fmt.Sprintf("medfri-%s", strings.Join(pathParts[2:3], "/"))
	address, portService, err := consul.GetServiceAddressAndPort(consulClient, serviceName)

	if err != nil {
		fmt.Println(fmt.Sprintf("%s no se encuentra en consulRegister", serviceName))
	}

	body, err, serviceStatusCode := service.GetServiceResponse(c,
		address,
		portService,
		serviceName,
		c.Request.Method,
		false)

	if len(body) == 0 {
		c.JSON(404, gin.H{"error": "api no encontrada"})
		c.Abort()
		return
	}

	if body["data"] != "data no avalible on the service" {

		if body["error"] == nil && ignoreCache != "Y" {
			service.PostServiceResponse(c, address, port, cacheServiceName, "POST", true, body)
		}

		fullUrl := strings.Split(c.Request.URL.String(), "/")

		estado := *serviceStatusCode
		respuesta := body["data"]
		microservicio := fullUrl[2]
		coleccion := fullUrl[3]
		endpoint := strings.Join(fullUrl[4:], "")

		fmt.Println("estado: ", estado)
		fmt.Println("respuesta: ", respuesta)
		fmt.Println("microservicio: ", microservicio)
		fmt.Println("colection: ", coleccion)
		fmt.Println("endpoint: ", endpoint)

		c.JSON(*serviceStatusCode, body)
		c.Abort()
		return
	}

	return
}
