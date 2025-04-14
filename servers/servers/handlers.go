package servers

import (
	"eth-toy-client/config"
	"net/http"
)

func SetupPingRoute(name config.ServerName, mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(name + " says pong"))
		if err != nil {
			return
		}
	})
}
