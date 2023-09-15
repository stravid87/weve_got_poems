package proxy

import (
	"fmt"
	"net/http"

	interfaces "go-load-balancer/interfaces"
)

type LoadBalancer struct {
	Port            string
	RoundRobinCount int
	Servers         []interfaces.Server
}

// getNextServerAddr returns the address of the next available server to send a request to, using a simple round-robin algorithm
func (lb *LoadBalancer) getNextAvailableServer() interfaces.Server {
	server := lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
	for !server.IsAlive() {
		lb.RoundRobinCount++
		server = lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
	}
	lb.RoundRobinCount++

	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()

	// could optionally log stuff about the request here!
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())

	// could delete pre-existing X-Forwarded-For header to prevent IP spoofing
	targetServer.Serve(rw, req)
}
