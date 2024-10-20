package main

import (
	"fmt"
	"getway-go/consul"
	"getway-go/httpServer"
	"getway-go/worker"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"runtime"
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

	consulClient := consul.ConnectToConsul()

	numCPUs := runtime.NumCPU()

	fmt.Printf("Detected %d CPUs, creating %d workers\n", numCPUs, numCPUs)

	taskQueue := make(chan *http.Request, 100)

	stop := make(chan struct{})

	worker.CreateWorkers(numCPUs, stop, taskQueue)

	go httpServer.InitHttpServer(taskQueue)

	worker.HandleShutdown(stop, consulClient)

}
