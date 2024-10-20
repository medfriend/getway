package main

import (
	"fmt"
	"getway-go/consul"
	"getway-go/httpServer"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func handleShutdown(stop chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Printf("Caught signal %s: shutting down gracefully...\n", sig)
	close(stop)
}

func createWorkers(numWorkers int, stop chan struct{}, taskQueue chan *http.Request) {
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			for {
				select {
				case req := <-taskQueue:
					fmt.Printf("Worker %d processing request for: %s\n", id, req.URL.Path)

				case <-stop:
					fmt.Printf("Worker %d shutting down\n", id)
					return
				}
			}
		}(i)
	}
}

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

	numCPUs := runtime.NumCPU()
	fmt.Printf("Detected %d CPUs, creating %d workers\n", numCPUs, numCPUs)

	taskQueue := make(chan *http.Request, 100)

	stop := make(chan struct{})

	createWorkers(numCPUs, stop, taskQueue)

	go httpServer.InitHttpServer(taskQueue)

	handleShutdown(stop)

	fmt.Println("Server stopped gracefully")
}
