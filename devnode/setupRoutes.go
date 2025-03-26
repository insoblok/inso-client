package devnode

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
)

func SetupRoutes(devAccount common.Address) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/dev-account", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Address string `json:"address"`
		}{
			Address: devAccount.Hex(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// ... more handlers

	return mux
}
