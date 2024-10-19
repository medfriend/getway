package httpServer

import (
	"fmt"
	"getway-go/jwt"
	"net/http"
)

// TODO crear logica de listas blancas
func InitHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hola desde %s!", "getway-go")

		corsConfig(w, r)

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

func corsConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
