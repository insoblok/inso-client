package devserver

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"eth-toy-client/core/consts"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/servers/servers"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

func handleDevAccounts(devAccount common.Address) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Address string `json:"address"`
		}{
			Address: devAccount.Hex(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleAccounts(accounts *map[string]*TestAccount) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func handleInfo(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			RPCURL        string `json:"rpcUrl"`
			RPCPort       string `json:"rpcPort"`
			AccountsCount int    `json:"accountsCount"`
		}{
			RPCURL:        "http://localhost:" + nodeClient.Config.Port,
			RPCPort:       nodeClient.Config.Port,
			AccountsCount: len(*accounts),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func signTxHandler(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req toytypes.SignTxRequest
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

		ctx := context.Background()
		var err error
		nonce := uint64(0)
		if req.Nonce != nil {
			nonce = *req.Nonce
		} else {
			nonce, err = nodeClient.Client.PendingNonceAt(ctx, from.Address)
			if err != nil {
				http.Error(w, "Failed to get nonce", http.StatusInternalServerError)
				return
			}
		}

		chainID := big.NewInt(consts.DefaultChainID)
		if req.ChainID != nil {
			chainID = big.NewInt(*req.ChainID)
		} else {
			chainID, err = nodeClient.Client.ChainID(ctx)
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
			Gas:       1_500_000,
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

func handleSendTx(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req toytypes.SignTxRequest
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

		ctx := context.Background()
		chainID, _ := nodeClient.Client.ChainID(ctx)
		nonce, _ := nodeClient.Client.PendingNonceAt(ctx, fromAcc.Address)

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
			Gas:       1_500_000,
			To:        &toAcc.Address,
			Value:     value,
		})

		signer := types.NewLondonSigner(chainID)
		signedTx, _ := types.SignTx(tx, signer, fromAcc.PrivKey)

		err := nodeClient.Client.SendTransaction(ctx, signedTx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send tx: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(SendTxResponse{TxHash: signedTx.Hash().Hex()})
	}
}

func handleSignTx(nodeClient *servers.NodeClient, accounts *map[string]*TestAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		tx, signedTx, err := BuildAndSignTx(from.PrivKey, from.Address, &to.Address, val, nodeClient.Config.Port, nil)
		if err != nil {
			log.Printf("❌ Signing failed: %v", err)
			httpapi.WriteError(w, http.StatusInternalServerError, "SigningFailed", err.Error())
			return
		}

		log.Printf("✅ Signed TX from %s → %s | value=%s | hash=%s",
			from.Address.Hex(), to.Address.Hex(), val.String(), tx.Hash().Hex())

		resp := &toytypes.SignTxAPIResponse{
			SignedTx: hex.EncodeToString(RlpEncodeBytes(signedTx)),
			TxHash:   tx.Hash().Hex(),
		}

		httpapi.WriteOK[toytypes.SignTxAPIResponse](w, resp)
	}
}

func handleRegisterAlias(reg *contract.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var meta contract.DeployedContractMetaJSON
		if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
			httpapi.WriteError(w, 400, "InvalidRequest", "Could not parse JSON")
			return
		}

		// 🧪 Validate required fields
		if meta.Alias == "" || meta.Address == "" || meta.TxHash == "" {
			httpapi.WriteError(w, 400, "MissingFields", "Alias, address, and txHash are required")
			return
		}

		// 🕒 Set timestamp if not provided
		if meta.Timestamp == 0 {
			meta.Timestamp = time.Now().Unix()
		}

		logutil.Infof("📦 Registering alias: %s → %s", meta.Alias, meta.Address)

		if err := reg.Add(meta); err != nil {
			httpapi.WriteError(w, 400, "DuplicateAlias", err.Error())
			return
		}

		httpapi.WriteOK(w, &toytypes.AliasRegisterResponse{
			Status: "ok",
			Alias:  meta.Alias,
		})
	}
}

func handleGetContracts(reg *contract.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := reg.All()
		summaries := []contract.DeployedContractMetaJSON{}

		for _, entry := range all {
			entry.ABI = ""      // Strip heavy fields
			entry.Bytecode = "" // Keep response lightweight
			summaries = append(summaries, entry)
		}

		httpapi.WriteOK(w, &summaries)
	}
}

func handleGetContractByAlias(reg *contract.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := strings.TrimPrefix(r.URL.Path, "/api/contracts/")
		if address == "" {
			httpapi.WriteError(w, 400, "MissingContract", "Contract address is required in the path")
			return
		}
		contractAddress := toytypes.ContractAddress{Address: address}
		meta, ok := reg.Get(contractAddress)
		if !ok {
			httpapi.WriteError(w, 404, "NotFound", fmt.Sprintf("Address '%s' not found", contractAddress.Address))
			return
		}

		httpapi.WriteOK(w, &meta)
	}
}
