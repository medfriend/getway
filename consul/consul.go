package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ConnectToConsul() {

	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		log.Fatalf("Error creating consul client: %v", err)
	}

	serviceID := os.Getenv("SERVICE_ID")
	serviceName := os.Getenv("SERVICE_NAME")
	serviceAddress := os.Getenv("SERVICE_ADDRESS")
	servicePort := os.Getenv("SERVICE_PORT")

	if serviceID == "" || serviceName == "" || serviceAddress == "" || servicePort == "" {
		log.Fatalf("Missing required environment variables")
	}

	//TODO implementar el healthcheck
	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    os.Getenv("SERVICE_NAME"),
		Address: os.Getenv("SERVICE_ADDRESS"),
		Port:    8080,
		//Check: &api.AgentServiceCheck{
		//	HTTP:     fmt.Sprintf("http://%s:%d/health", serviceAddress, 8080),
		//	Interval: "10s",
		//	Timeout:  "5s",
		//},
	}

	err = client.Agent().ServiceRegister(service)

	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}

	fmt.Println("Service registered successfully")

	//handleShutdown(client, serviceID)
}

func handleShutdown(client *api.Client, serviceID string) {

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigChan
		fmt.Printf("Caught signal %s: shutting down service\n", sig)

		err := client.Agent().ServiceDeregister(serviceID)
		if err != nil {
			log.Fatalf("Error deregistering service: %v", err)
		}

		fmt.Println("Service deregistered successfully")

		os.Exit(0)
	}()
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is healthy"))
}

func GetServiceAddressAndPort(serviceName string) (string, int, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", 0, fmt.Errorf("error creating Consul client: %v", err)
	}

	services, _, err := client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return "", 0, fmt.Errorf("error retrieving service: %v", err)
	}

	if len(services) == 0 {
		return "", 0, fmt.Errorf("no instances found for service: %s", serviceName)
	}

	service := services[0]
	return service.Address, service.ServicePort, nil
}
