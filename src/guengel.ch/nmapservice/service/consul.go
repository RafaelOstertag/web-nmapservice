package service

import (
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
)

const (
	serviceName = "nmap"
)

func getConsulAgentAddress() string {
	consul := os.Getenv("CONSUL")
	if consul == "" {
		return "gizmo.kruemel.home:8500"
	}
	return consul
}

func makeCheck(host string, port int) *api.AgentServiceCheck {
	check := new(api.AgentServiceCheck)
	check.Name = "Check " + serviceName
	check.GRPC = host + ":" + strconv.Itoa(port)
	check.Interval = "15s"
	check.DeregisterCriticalServiceAfter = "30s"

	return check
}

func connectConsul(consulAgentAddr string) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = consulAgentAddr
	consul, err := api.NewClient(config)
	if err != nil {
		log.Printf("Error creating consul client: %v", err)
		return nil, err
	}

	return consul, nil
}

// Register registers the service with consul
func Register(host string, port int) error {
	registration := new(api.AgentServiceRegistration)

	registration.Address = host
	registration.Port = port
	registration.Name = serviceName
	registration.Check = makeCheck(host, port)

	consulAddress := getConsulAgentAddress()
	log.Printf("Connect to consul %s", consulAddress)
	consul, err := connectConsul(consulAddress)
	if err != nil {
		return err
	}

	consulAgent := consul.Agent()
	if err := consulAgent.ServiceRegister(registration); err != nil {
		log.Printf("Error registering service: %v", err)
		return err
	}

	log.Printf("Registered service with consul as '%s:%d'", host, port)

	return nil
}
