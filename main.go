package main

import (
	"fmt"
	"getway-go/httpServer"
	"getway-go/util/consul"
	"getway-go/util/env"
	"getway-go/util/worker"
	"net/http"
	"runtime"
)

func main() {

	env.LoadEnv()

	consulClient := consul.ConnectToConsul()

	numCPUs := runtime.NumCPU()

	fmt.Printf("Detected %d CPUs, creating %d workers\n", numCPUs, numCPUs)

	taskQueue := make(chan *http.Request, 100)

	stop := make(chan struct{})

	worker.CreateWorkers(numCPUs, stop, taskQueue)

	go httpServer.InitHttpServer(taskQueue)

	worker.HandleShutdown(stop, consulClient)

}
