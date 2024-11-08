package redirectgetway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"getway-go/util"
	"github.com/gin-gonic/gin"
	"github.com/medfriend/shared-commons-go/util/consul"
	"io"
	"net/http"
	"os"
	"strings"
)

func Redirectgetway(c *gin.Context) {

	cacheServiceName := os.Getenv("SERVICE_CACHE")

	address, port, err := consul.GetServiceAddressAndPort(cacheServiceName)

	if err != nil {
		fmt.Println("cache no se encuentra en consulRegister")
	}

	if err == nil {
		body, errCache, cacheStatusCode := getServiceResponse(c, address, port, cacheServiceName, "GET", true)

		if body["data"] != "data no avalible on the cache" {
			c.JSON(*cacheStatusCode, body)
			c.Abort()
			return
		}

		if errCache != nil || body["data"] == "data no avalible on the cache" {

			pathParts := strings.Split(c.Request.URL.Path, "/")

			serviceName := fmt.Sprintf("medfri-%s", strings.Join(pathParts[2:3], "/"))

			addressService, portService, err := consul.GetServiceAddressAndPort(serviceName)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s no se encuentra en consulRegister", serviceName))
			}

			body, err, serviceStatusCode := getServiceResponse(c,
				addressService,
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
				postServiceResponse(c, address, port, cacheServiceName, "POST", true, body)
				c.JSON(*serviceStatusCode, body)
				c.Abort()
				return
			}

			return
		}
	}

}

func getServiceResponse(c *gin.Context, address string, port int, serviceGetway string, method string, isCache bool) (map[string]interface{}, error, *int) {

	// Leer el body de la solicitud original
	bodyBytes, err := io.ReadAll(c.Request.Body)

	targetURL := util.CreatePath(address, port, serviceGetway, c, isCache)

	// Restaurar el body de la solicitud original para futuras lecturas (opcional)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Crear la nueva solicitud con el body leído
	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("Error creando la nueva solicitud"), nil
	}

	// Copiar los headers de la solicitud original
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Enviar la solicitud al servicio secundario
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error al realizar la solicitud al servicio"), nil
	}
	defer resp.Body.Close()

	// Leer el body de la respuesta del servicio
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error al leer la respuesta del servicio de cache"), nil
	}

	statusCode := resp.StatusCode

	// Decodificar la respuesta JSON
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("Error al decodificar la respuesta del servicio de cache"), nil
	}

	return result, nil, &statusCode
}

func postServiceResponse(c *gin.Context, address string, port int, serviceGetway string, method string, isCache bool, body map[string]interface{}) (map[string]interface{}, error, *int) {
	// Convertir el body (map) en JSON
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Error al serializar el cuerpo de la solicitud: %v", err), nil
	}

	// Crear la URL de destino
	targetURL := util.CreatePath(address, port, serviceGetway, c, isCache)

	// Crear la nueva solicitud con el body serializado a JSON
	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("Error creando la nueva solicitud"), nil
	}

	// Copiar los headers de la solicitud original
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Enviar la solicitud al servicio secundario
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error al realizar la solicitud al servicio"), nil
	}
	defer resp.Body.Close()

	// Leer el body de la respuesta del servicio
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error al leer la respuesta del servicio"), nil
	}

	// Obtener el código de respuesta HTTP
	statusCode := resp.StatusCode

	// Decodificar la respuesta JSON
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("Error al decodificar la respuesta del servicio"), nil
	}

	return result, nil, &statusCode
}
