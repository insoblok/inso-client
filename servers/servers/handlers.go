package servers

import "net/http"

func SetupPingRoute(name string, mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(name + " says pong"))
		if err != nil {
			return
		}
	})
}
