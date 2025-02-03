package main

import (
	"fmt"
	"getway-go/consulRegister"
	"getway-go/httpServer"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/medfriend/shared-commons-go/util/env"
	"github.com/medfriend/shared-commons-go/util/worker"
	"net/http"
	"os"
	"runtime"
)

func main() {

	env.LoadEnv()

	deployment := os.Getenv("DEPLOYMENT")

	consulHost := "localhost"

	if deployment == "DOCKER" {
		consulHost = os.Getenv("CONSUL_ADDRESS_DOCKER")
	}

	if deployment == "KUBE" {
		consulHost = os.Getenv("CONSUL_ADDRESS_KUBE")
	}

	consulAddress := fmt.Sprintf("http://%s:8500", consulHost)

	consulClient := consul.ConnectToConsul(consulAddress)

	consulRegister.RegisterConstants(consulClient)

	numCPUs := runtime.NumCPU()

	fmt.Printf("Detected %d CPUs, creating %d workers\n", numCPUs, numCPUs)

	taskQueue := make(chan *http.Request, 100)

	stop := make(chan struct{})

	worker.CreateWorkers(numCPUs, stop, taskQueue)

	go httpServer.InitHttpServer(taskQueue)

	worker.HandleShutdown(stop, consulClient)
}
