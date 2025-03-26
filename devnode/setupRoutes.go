package devnode

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"net/http"
)

type accountResponse struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func SetupRoutes(devAccount common.Address, accounts *map[string]*TestAccount) *http.ServeMux {
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

	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		var list []accountResponse
		for name, acc := range *accounts {
			privKeyBytes := crypto.FromECDSA(acc.PrivKey) // import "github.com/ethereum/go-ethereum/crypto"
			list = append(list, accountResponse{
				Name:       name,
				Address:    acc.Address.Hex(),
				PrivateKey: hex.EncodeToString(privKeyBytes),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})

	return mux
}
