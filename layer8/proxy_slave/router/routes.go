package routes

import (
	"net/http"

	Ctl "proxy_slave/api/controller"
)

func RegisterRoutes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set up route for API
		switch r.URL.Path {

		case "/api/v1/ping":
			Ctl.Ping(w, r)

		case "/api/v1/register-user":
			Ctl.RegisterUserHandler(w, r)

		case "/api/v1/login-precheck":
			Ctl.LoginPrecheckHandler(w, r)

		case "/api/v1/login-user":
			Ctl.LoginUserHandler(w, r)
		case "/api/v1/get-content":
			Ctl.GetContentHandler(w, r)
		case "/api/v1/initialize-ecdh-tunnel":
			Ctl.InitializeECDHTunnelHandler(w, r)

		default:
			// Return a 404 Not Found error for unknown routes
			http.NotFound(w, r)
		}
	}
}
