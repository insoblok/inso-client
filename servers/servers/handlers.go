package servers

import "net/http"

func SetupPingRoute(mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			return
		}
	})
}
