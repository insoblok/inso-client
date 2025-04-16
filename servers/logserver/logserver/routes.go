package logserver

import (
	"encoding/json"
	"eth-toy-client/config"
	"eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/servers/servers"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"net/http"
	"strings"
	"time"
)

func SetupRoutes(config config.ServerConfig, contractRegistry *contract.Registry) *http.ServeMux {
	mux := http.NewServeMux()
	servers.SetupPingRoute(config.Name, mux)
	mux.HandleFunc("/api/register-contract", registerContract(contractRegistry))
	mux.Handle("/api/contract/", http.StripPrefix("/contract", getContract(contractRegistry)))
	return mux
}

func getContract(registry *contract.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Path[len("/"):]
		if address == "" {
			httpapi.WriteError(w, 400, "❌ MissingAlias", "Alias must be specified in the path")
			return
		}
		contactAddress := toytypes.ContractAddress{Address: address}
		meta, err := registry.Get(contactAddress)
		if !err {
			httpapi.WriteError(w, 404, "❌ AliasNotFound", "Failed to get contract address "+contactAddress.Address)
			return
		}
		httpapi.WriteOK(w, &meta)
	}
}

func registerContract(reg *contract.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var meta contract.DeployedContractMetaJSON
		if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
			httpapi.WriteError(w, 400, "❌ InvalidRequest", "Could not parse JSON")
			return
		}
		if meta.Alias == "" || meta.Address == "" {
			httpapi.WriteError(w, 400, "❌ MissingFields", "Alias, address are required")
			return
		}
		if meta.Timestamp == 0 {
			meta.Timestamp = time.Now().Unix()
		}

		parsedABI, err := abi.JSON(strings.NewReader(meta.ABI))

		if err != nil {
			logutil.Errorf("❌ Error parsing ABI: %v", err)
			httpapi.WriteError(w, 400, "❌ InvalidABI", "Could not parse ABI")
			return
		}

		info := contract.DeployedContractInfo{
			Address:   toytypes.ContractAddress{Address: meta.Address},
			Pending:   true,
			Alias:     meta.Alias,
			ABI:       meta.ABI,
			ParsedABI: &parsedABI,
		}

		logutil.Infof("📦 Registering alias: %s → %s", meta.Alias, meta.Address)
		if err := reg.Add(info); err != nil {
			httpapi.WriteError(w, 400, "DuplicateAlias", err.Error())
			return
		}

		res := toytypes.AliasRegisterResponse{
			Status: "ok",
			Alias:  meta.Alias,
		}
		logutil.Infof("✅ sending: %v", res)
		httpapi.WriteOK(w, &res)
	}
}
