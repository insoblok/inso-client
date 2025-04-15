package devserver

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/servers/servers"
	"fmt"
	"log"
	"net/http"
)

func deployContract(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Headers: %+v\n", r.Header)
		if r.Method != http.MethodPost {
			log.Printf("⚠️ Invalid method: %s", r.Method)
			httpapi.WriteError(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Only POST is allowed")
			return
		}

		var req toytypes.DeployContractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("❌ Failed to decode JSON: %v", err)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidRequest", "Invalid JSON payload")
			return
		}

		from, ok := (*accounts)[req.From]
		if !ok {
			log.Printf("⚠️ Sender not found: %s", req.From)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Sender '%s' not found", req.From))
			return
		}

		rawByte := []byte(req.Data)
		if (len(rawByte) % 2) == 1 {
			rawByte = append([]byte("0"), rawByte...)
		}

		destHexByte := make([]byte, len(rawByte)/2)
		hex.Decode(destHexByte, rawByte)
		data := destHexByte
		logutil.Infof("Hex Bytes Length: %d", len(data))

		_, contractAddress, signedTx, err := SignContract(from.PrivKey, from.Address, req.Nonce, nodeClient.Config.Port, data)
		if err != nil {
			log.Printf("❌ Signing failed: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "SigningFailed", err.Error())
			return
		}

		err = nodeClient.Client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Printf("❌ Failed to send tx: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "SendTxFailed", err.Error())
			return
		}

		log.Printf("✅ Sent TX: %s", signedTx.Hash().Hex())

		httpapi.WriteOK[toytypes.ContractDeploymentResponse](w, &toytypes.ContractDeploymentResponse{
			TxHash:                  signedTx.Hash().Hex(),
			ExpectedContractAddress: contractAddress.Hex(),
		})
	}
}
