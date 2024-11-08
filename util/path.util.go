package util

import (
	"fmt"
	"getway-go/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func CreatePath(address string, port int, serviceGetway string, c *gin.Context, isCache bool) string {

	// Obtener partes del path y ensamblar el servicio
	pathParts := strings.Split(c.Request.URL.Path, "/")
	service := strings.Join(pathParts[3:], "/")

	// Verificar si es cache y si el servicio es "auth"

	if isCache && service != "auth" {
		authorization := c.Request.Header.Get("Authorization")
		decode, _ := jwt.DecodeJWT(authorization)
		service = service + "/" + strconv.Itoa(decode.User.Usuario)
	}

	if isCache && service == "auth" {
		usuario := c.Request.Header.Get("usuario")
		fmt.Println(usuario)
		service = service + "/" + usuario
	}

	// Construir y devolver la URL de destino
	targetURL := fmt.Sprintf("http://%s:%d/%s/%s/", address, port, serviceGetway, service)

	fmt.Println(targetURL)

	return targetURL
}
