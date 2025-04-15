package logserver

import (
	"encoding/json"
	"eth-toy-client/config"
	"eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/servers/servers"
	"net/http"
	"time"
)

func SetupRoutes(config config.ServerConfig, contractRegistry *contract.ContractRegistry) *http.ServeMux {
	mux := http.NewServeMux()
	servers.SetupPingRoute(config.Name, mux)
	mux.HandleFunc("/api/register-contract", registerContract(contractRegistry))
	mux.Handle("/api/contract/", http.StripPrefix("/contract", getContract(contractRegistry)))
	return mux
}

func getContract(registry *contract.ContractRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.URL.Path[len("/"):]
		if alias == "" {
			httpapi.WriteError(w, 400, "❌ MissingAlias", "Alias must be specified in the path")
			return
		}
		reg, err := registry.Get(alias)
		if !err {
			httpapi.WriteError(w, 404, "❌ AliasNotFound", "Failed to get alias "+alias)
			return
		}
		httpapi.WriteOK(w, &reg)
	}
}

func registerContract(reg *contract.ContractRegistry) http.HandlerFunc {
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

		logutil.Infof("📦 Registering alias: %s → %s", meta.Alias, meta.Address)
		if err := reg.Add(meta); err != nil {
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
