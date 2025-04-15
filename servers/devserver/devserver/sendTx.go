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
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"net/http"
)

func handleSendTxAPI(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Headers: %+v\n", r.Header)
		if r.Method != http.MethodPost {
			log.Printf("⚠️ Invalid method: %s", r.Method)
			httpapi.WriteError(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Only POST is allowed")
			return
		}

		var req toytypes.SignTxRequest
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

		val := new(big.Int)
		if req.Value == "" {
			val.SetInt64(0)
		} else if _, ok := val.SetString(req.Value, 10); !ok {
			log.Printf("❌ Invalid value format: %s", req.Value)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidValue", "Invalid value format")
			return
		}

		var (
			toAddr *common.Address
			data   []byte
		)

		if req.To == "" {
			// 🚀 Contract Deployment
			log.Printf("📨 /send-tx: deploying contract from=%s", req.From)

			if req.Data == "" {
				log.Printf("❌ Missing contract bytecode in data field")
				httpapi.WriteError(w, http.StatusBadRequest, "MissingData", "Contract deployment requires 'data' field")
				return
			}

			toAddr = nil
			rawByte := []byte(req.Data)
			if (len(rawByte) % 2) == 1 {
				rawByte = append([]byte("0"), rawByte...)
			}

			destHexByte := make([]byte, len(rawByte)/2)
			hex.Decode(destHexByte, rawByte)
			data = destHexByte
			logutil.Infof("Hex Bytes Length: %d", len(data))

		} else {
			// 🔁 Normal Transfer
			toAccount, ok := (*accounts)[req.To]
			if !ok {
				log.Printf("⚠️ Recipient not found: %s", req.To)
				httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Recipient '%s' not found", req.To))
				return
			}
			addr := toAccount.Address
			toAddr = &addr
			log.Printf("📨 /send-tx: from=%s → to=%s | value=%s", req.From, req.To, req.Value)
		}

		_, signedTx, err := BuildAndSignTx(from.PrivKey, from.Address, toAddr, val, nodeClient.Config.Port, data)
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

		httpapi.WriteOK[toytypes.SendTxAPIResponse](w, &toytypes.SendTxAPIResponse{
			TxHash: signedTx.Hash().Hex(),
		})
	}
}
