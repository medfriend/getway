package httpServer

import (
	"fmt"
	"getway-go/jwt"
	"net/http"
)

func InitHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hola desde %s!", "getway-go")

		authHeader := r.Header.Get("Authorization")

		token, err := jwt.ValidateJWT(authHeader)

		//TODO empezar a crear request error handler
		if err != nil || !token.Valid {
			fmt.Println("Token inválido:", err)
			return
		}

		fmt.Println(token)

		fmt.Println("test done")

	})
}
