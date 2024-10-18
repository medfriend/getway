package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
)

func ValidateJWT(tokenString string) (*jwt.Token, error) {

	jwtKey := []byte(os.Getenv("JWT_KEY"))

	tokenSplited, err := splitToken(tokenString)

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenSplited, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func splitToken(tokenString string) (string, error) {
	parts := strings.Split(tokenString, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		fmt.Println("token invalido")
		return "", errors.New("token invalido")
	}

	return parts[1], nil
}
