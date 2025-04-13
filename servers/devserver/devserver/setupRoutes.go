package devserver

import (
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/servers/servers"
	"eth-toy-client/swagger"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
)

type accountResponse struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func SetupRoutes(
	reg *contract.ContractRegistry,
	devAccount common.Address,
	nodeClient *servers.NodeClient,
	accounts *map[string]*TestAccount) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/dev-account", handleDevAccounts(devAccount))
	mux.HandleFunc("/accounts", handleAccounts(accounts))
	mux.HandleFunc("/info", handleInfo(nodeClient, accounts))
	mux.HandleFunc("/sign-tx", signTxHandler(nodeClient, accounts))
	mux.HandleFunc("/send-tx", handleSendTx(nodeClient, accounts))
	mux.HandleFunc("/api/sign-tx", handleSignTx(nodeClient, accounts))
	mux.HandleFunc("/api/send-tx", handleSendTxAPI(nodeClient, accounts))
	mux.HandleFunc("/api/register-alias", handleRegisterAlias(reg))
	mux.HandleFunc("/api/contracts", handleGetContracts(reg))
	mux.HandleFunc("/api/contracts/", handleGetContractByAlias(reg))
	mux.HandleFunc("/swagger/", swagger.HandleSwagger)

	return mux
}
