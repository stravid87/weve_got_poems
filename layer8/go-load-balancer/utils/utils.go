package utils

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"

	interfaces "go-load-balancer/interfaces"
	"go-load-balancer/proxy"
	"go-load-balancer/server"
)

func NewSimpleServer(addr string) *server.SimpleServer {
	serverUrl, err := url.Parse(addr)
	fmt.Println(serverUrl)
	handleErr(err)

	return &server.SimpleServer{
		Addr:  addr,
		Proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func NewLoadBalancer(port string, servers []interfaces.Server) *proxy.LoadBalancer {
	return &proxy.LoadBalancer{
		Port:            port,
		RoundRobinCount: 0,
		Servers:         servers,
	}
}

// handleErr prints the error and exits the program
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
