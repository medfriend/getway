package main

import (
	"getway-go/consul"
	"getway-go/httpServer"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// TODO agregar el modulo de lint y validacion test
func main() {

	loadEnv()

	consul.ConnectToConsul()

	http.HandleFunc("/health", consul.HealthCheck)
	httpServer.InitHttpServer()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
