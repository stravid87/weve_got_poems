package main

import (
	"fmt"
	"net/http"
	"os"

	interfaces "go-load-balancer/interfaces"
	"go-load-balancer/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load the environment variables from .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	// Get the port to listen on from .env
	loadBalancerPort := os.Getenv("LOAD_BALANCER_PORT")
	// create a slice of layer8-slaves to forward requests to
	servers := []interfaces.Server{
		utils.NewSimpleServer("http://localhost:8001"),
		// utils.NewSimpleServer("http://localhost:8002"),
		// utils.NewSimpleServer("http://localhost:8003"),
	}
	// create a load balancer and serve proxy requests
	lb := utils.NewLoadBalancer(loadBalancerPort, servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}
	// register a proxy handler to handle all requests
	http.HandleFunc("/", handleRedirect)
	// start the server on the specified port
	fmt.Printf("serving requests at 'localhost:%s'\n", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}
