package devnode

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"eth-toy-client/core/consts"
	"eth-toy-client/core/httpapi"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"net/http"
)

type accountResponse struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func SetupRoutes(devAccount common.Address, rpcPort string, accounts *map[string]*TestAccount) *http.ServeMux {
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

	// 🆕 /info endpoint
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			RPCURL        string `json:"rpcUrl"`
			RPCPort       string `json:"rpcPort"`
			AccountsCount int    `json:"accountsCount"`
		}{
			RPCURL:        "http://localhost:" + rpcPort,
			RPCPort:       rpcPort,
			AccountsCount: len(*accounts),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/sign-tx", signTxHandler(rpcPort, accounts))
	mux.HandleFunc("/send-tx", handleSendTx(rpcPort, accounts))

	mux.HandleFunc("/api/sign-tx", handleSignTx(rpcPort, accounts))
	mux.HandleFunc("/api/send-tx", handleSendTxAPI(rpcPort, accounts))

	return mux
}

func signTxHandler(rpcPort string, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignTxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		from, ok := (*accounts)[req.From]
		if !ok {
			http.Error(w, "Unknown sender account", http.StatusBadRequest)
			return
		}

		toAddr := common.HexToAddress(req.To)
		value := big.NewInt(consts.DefaultTransferAmount) // default
		if req.Value != "" {
			v, ok := new(big.Int).SetString(req.Value, 10)
			if !ok {
				http.Error(w, "Invalid value field", http.StatusBadRequest)
				return
			}
			value = v
		}

		// 🧠 Connect to RPC to get nonce / chain ID
		rpcClient, err := rpc.Dial("http://localhost:" + rpcPort)
		if err != nil {
			http.Error(w, "RPC dial failed", http.StatusInternalServerError)
			return
		}
		ethClient := ethclient.NewClient(rpcClient)
		defer ethClient.Close()

		ctx := context.Background()
		nonce := uint64(0)
		if req.Nonce != nil {
			nonce = *req.Nonce
		} else {
			nonce, err = ethClient.PendingNonceAt(ctx, from.Address)
			if err != nil {
				http.Error(w, "Failed to get nonce", http.StatusInternalServerError)
				return
			}
		}

		chainID := big.NewInt(consts.DefaultChainID)
		if req.ChainID != nil {
			chainID = big.NewInt(*req.ChainID)
		} else {
			chainID, err = ethClient.ChainID(ctx)
			if err != nil {
				http.Error(w, "Failed to get chain ID", http.StatusInternalServerError)
				return
			}
		}

		// 🧾 Construct tx
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: big.NewInt(1),
			GasFeeCap: big.NewInt(1e9),
			Gas:       21_000,
			To:        &toAddr,
			Value:     value,
		})

		// ✍️ Sign it
		signer := types.NewLondonSigner(chainID)
		signedTx, err := types.SignTx(tx, signer, from.PrivKey)
		if err != nil {
			http.Error(w, "Failed to sign tx", http.StatusInternalServerError)
			return
		}

		var buff bytes.Buffer
		signedTx.EncodeRLP(&buff)

		resp := SignTxResponse{
			Tx: "0x" + hex.EncodeToString(buff.Bytes()),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}

}

func handleSendTx(rpcPort string, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignTxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		fromAcc, ok := (*accounts)[req.From]
		if !ok {
			http.Error(w, "From account not found", http.StatusNotFound)
			return
		}

		toAcc, ok := (*accounts)[req.To]
		if !ok {
			http.Error(w, "To account not found", http.StatusNotFound)
			return
		}

		client, err := ethclient.Dial("http://localhost:" + rpcPort)
		if err != nil {
			http.Error(w, "Failed to connect to dev node", http.StatusInternalServerError)
			return
		}
		defer client.Close()

		ctx := context.Background()
		chainID, _ := client.ChainID(ctx)
		nonce, _ := client.PendingNonceAt(ctx, fromAcc.Address)

		value := big.NewInt(consts.DefaultTransferAmount) // default
		if req.Value != "" {
			v, ok := new(big.Int).SetString(req.Value, 10)
			if !ok {
				http.Error(w, "Invalid value field", http.StatusBadRequest)
				return
			}
			value = v
		}

		// 🧾 Construct tx
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: big.NewInt(1),
			GasFeeCap: big.NewInt(1e9),
			Gas:       21_000,
			To:        &toAcc.Address,
			Value:     value,
		})

		signer := types.NewLondonSigner(chainID)
		signedTx, _ := types.SignTx(tx, signer, fromAcc.PrivKey)

		err = client.SendTransaction(ctx, signedTx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send tx: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(SendTxResponse{TxHash: signedTx.Hash().Hex()})
	}
}

func handleSignTx(rpcPort string, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("⚠️ Invalid method: %s", r.Method)
			httpapi.WriteError(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Only POST is allowed")
			return
		}

		var req SignTxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("❌ Failed to decode JSON: %v", err)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidRequest", "Invalid JSON payload")
			return
		}

		log.Printf("📨 /sign-tx: from=%s → to=%s | value=%s", req.From, req.To, req.Value)

		from, ok := (*accounts)[req.From]
		if !ok {
			log.Printf("⚠️ Sender not found: %s", req.From)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Sender '%s' not found", req.From))
			return
		}

		to, ok := (*accounts)[req.To]
		if !ok {
			log.Printf("⚠️ Recipient not found: %s", req.To)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Recipient '%s' not found", req.To))
			return
		}

		val := new(big.Int)
		_, ok = val.SetString(req.Value, 10)
		if !ok {
			log.Printf("❌ Invalid value format: %s", req.Value)
			httpapi.WriteError(w, http.StatusBadRequest, "INVALID_VALUE", "Invalid value format")
			return
		}

		tx, signedTx, err := BuildAndSignTx(from.PrivKey, from.Address, to.Address, val, rpcPort)
		if err != nil {
			log.Printf("❌ Signing failed: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "SigningFailed", err.Error())
			return
		}

		log.Printf("✅ Signed TX from %s → %s | value=%s | hash=%s",
			from.Address.Hex(), to.Address.Hex(), val.String(), tx.Hash().Hex())

		resp := &SignTxAPIResponse{
			SignedTx: hex.EncodeToString(RlpEncodeBytes(signedTx)),
			TxHash:   tx.Hash().Hex(),
		}

		httpapi.WriteOK[SignTxAPIResponse](w, resp)
	}
}

func handleSendTxAPI(rpcPort string, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("⚠️ Invalid method: %s", r.Method)
			httpapi.WriteError(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Only POST is allowed")
			return
		}

		var req SignTxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("❌ Failed to decode JSON: %v", err)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidRequest", "Invalid JSON payload")
			return
		}

		log.Printf("📨 /send-tx: from=%s → to=%s | value=%s", req.From, req.To, req.Value)

		from, ok := (*accounts)[req.From]
		if !ok {
			log.Printf("⚠️ Sender not found: %s", req.From)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Sender '%s' not found", req.From))
			return
		}

		to, ok := (*accounts)[req.To]
		if !ok {
			log.Printf("⚠️ Recipient not found: %s", req.To)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidAccount", fmt.Sprintf("Recipient '%s' not found", req.To))
			return
		}

		val := new(big.Int)
		if _, ok := val.SetString(req.Value, 10); !ok {
			log.Printf("❌ Invalid value format: %s", req.Value)
			httpapi.WriteError(w, http.StatusBadRequest, "InvalidValue", "Invalid value format")
			return
		}

		tx, signedTx, err := BuildAndSignTx(from.PrivKey, from.Address, to.Address, val, rpcPort)
		if err != nil {
			log.Printf("❌ Signing failed: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "SigningFailed", err.Error())
			return
		}

		client, err := ethclient.Dial("http://localhost:" + rpcPort)
		if err != nil {
			log.Printf("❌ Failed to connect to dev node: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "ConnectionFailed", "Could not connect to dev node")
			return
		}
		defer client.Close()

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Printf("❌ Failed to send tx: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "SendTxFailed", err.Error())
			return
		}

		log.Printf("✅ Sent TX: %s", tx.Hash().Hex())

		resp := &SendTxAPIResponse{
			TxHash: tx.Hash().Hex(),
		}
		httpapi.WriteOK[SendTxAPIResponse](w, resp)
	}
}
