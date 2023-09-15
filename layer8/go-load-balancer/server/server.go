package server

import (
	"net/http"
	"net/http/httputil"
)

type SimpleServer struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func (s *SimpleServer) Address() string { return s.Addr }

func (s *SimpleServer) IsAlive() bool { return true }

func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(rw, req)
}
