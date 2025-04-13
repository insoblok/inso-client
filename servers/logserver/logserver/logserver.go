package logserver

import "net/http"
import "eth-toy-client/servers/servers"

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	servers.SetupPingRoute(mux)
	mux.HandleFunc("/register-contract", registerContract())
	return mux
}

func registerContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
