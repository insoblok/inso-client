package devserver

import (
	"context"
	"encoding/json"
	"eth-toy-client/core/httpapi"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/servers/servers"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"net/http"
)

func handlePendingNonce(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("⚠️ Invalid method: %s", r.Method)
			httpapi.WriteError(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Only POST is allowed")
			return
		}

		var req toytypes.PendingNonceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("❌ Failed to decode JSON: %v", err)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidRequest", "Invalid JSON payload")
			return
		}

		from, ok := (*accounts)[req.Alias]
		if !ok {
			log.Printf("⚠️ Sender not found: %s", req.Alias)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Sender '%s' not found", req.Alias))
			return
		}

		nonce, err := nodeClient.Client.PendingNonceAt(context.Background(), from.Address)
		if err != nil {
			log.Printf("❌ Failed request pending Nonce: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "PendingNonceAt", err.Error())
			return
		}

		//nonceStr := fmt.Sprintf("%d", nonce)
		address := crypto.CreateAddress(from.Address, nonce)
		response := toytypes.PendingNonceResponse{
			Nonce:   &nonce,
			Address: address.Hex(),
		}
		log.Printf("Sending nonce: %v", response)
		httpapi.WriteOK[toytypes.PendingNonceResponse](w, &response)
	}
}
