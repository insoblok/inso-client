package devserver

import (
	contract "eth-toy-client/core/contracts"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
)

type accountResponse struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func SetupRoutes(reg *contract.ContractRegistry, devAccount common.Address, rpcPort string, accounts *map[string]*TestAccount) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/dev-account", handleDevAccounts(devAccount))
	mux.HandleFunc("/accounts", handleAccounts(accounts))
	mux.HandleFunc("/info", handleInfo(rpcPort, accounts))
	mux.HandleFunc("/sign-tx", signTxHandler(rpcPort, accounts))
	mux.HandleFunc("/send-tx", handleSendTx(rpcPort, accounts))
	mux.HandleFunc("/api/sign-tx", handleSignTx(rpcPort, accounts))
	mux.HandleFunc("/api/send-tx", handleSendTxAPI(rpcPort, accounts))
	mux.HandleFunc("/api/register-alias", handleRegisterAlias(reg))
	mux.HandleFunc("/api/contracts", handleGetContracts(reg))
	mux.HandleFunc("/api/contracts/", handleGetContractByAlias(reg))

	return mux
}
