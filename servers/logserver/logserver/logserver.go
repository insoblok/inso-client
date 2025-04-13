package logserver

import (
	"eth-toy-client/config"
	"net/http"
)
import "eth-toy-client/servers/servers"

func SetupRoutes(config config.ServerConfig) *http.ServeMux {
	mux := http.NewServeMux()
	servers.SetupPingRoute(config.Name, mux)
	mux.HandleFunc("/register-contract", registerContract())
	return mux
}

func registerContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
