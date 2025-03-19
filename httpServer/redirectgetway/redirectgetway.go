package redirectgetway

import (
	"encoding/json"
	"fmt"
	"getway-go/dto"
	"getway-go/httpServer/service"
	"getway-go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/medfriend/shared-commons-go/util/global"
	"github.com/medfriend/shared-commons-go/util/rabbitmq"
	"os"
	"strconv"
	"strings"
	"time"
)

func Redirectgetway(c *gin.Context, consulClient *api.Client, start time.Time) {

	ignoreCache := c.Request.Header.Get("ignore-cache")

	cacheServiceName := os.Getenv("SERVICE_CACHE")

	address, port, err := consul.GetServiceAddressAndPort(consulClient, cacheServiceName)

	if ignoreCache == "Y" {
		registerOnService(c, address, port, cacheServiceName, consulClient, start)
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
			registerOnService(c, address, port, cacheServiceName, consulClient, start)
		}
	}

}

func registerOnService(c *gin.Context, address string, port int, cacheServiceName string, consulClient *api.Client, start time.Time) {

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

		token := c.GetHeader("Authorization")
		respuestaJson, err := json.Marshal(body["data"])

		if err != nil {
			fmt.Errorf("error al parsear json de la respuesta del microservicio")
		}

		var usuario string

		if token == "" {
			decodeToken, _ := jwt.DecodeJWT("Bearer " + body["data"].(string))
			usuario = strconv.Itoa(decodeToken.User.UsuarioID)
		} else {
			decodeToken, _ := jwt.DecodeJWT(token)
			usuario = strconv.Itoa(decodeToken.User.UsuarioID)
		}

		var errorString string

		if body["error"] == nil {
			errorString = "no error"
		} else {
			errorString = body["error"].(string)
		}

		collection_id, err := json.Marshal(body["collectio_id"])

		trazaMessage := dto.TrazaDTO{
			UsuarioID:     usuario,
			Accion:        fullUrl[4],
			Estado:        *serviceStatusCode,
			Endpoint:      c.Request.URL.String(),
			Payload:       string(respuestaJson),
			Ip:            c.ClientIP(),
			Duracion:      time.Since(start).String(),
			Error:         errorString,
			Coleccion:     fullUrl[3],
			Microservicio: fullUrl[2],
			CollectionId:  string(collection_id),
		}

		trazaMessageJson, err := json.Marshal(trazaMessage)

		if err != nil {
			fmt.Errorf("error al marshal el trazamessage")
		}
		fmt.Println(global.GetRabbitConn())
		rabbit := rabbitmq.GetInstance(global.GetRabbitConn())

		rabbit.SendMessage(
			"trazaacciones",
			string(trazaMessageJson),
			global.GetRabbitConn(),
		)

		c.JSON(*serviceStatusCode, body)
		c.Abort()
		return
	}

	return
}
