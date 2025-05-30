package servers

import (
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	toytypes "eth-toy-client/core/types"
)

func RegisterContract(payload contract.DeployedContractMetaJSON) (*toytypes.AliasRegisterResponse, *httpapi.APIError, error) {
	registerURL := config.Servers.LogServer.GetServerConfig().GetServerUrl("register-contract")
	return httpapi.PostWithAPIResponse[toytypes.AliasRegisterResponse](registerURL, payload)
}

func GetContract(contractAddress toytypes.ContractAddress) (*contract.DeployedContractMetaJSON, *httpapi.APIError, error) {
	contractURL := config.Servers.LogServer.GetServerConfig().GetServerUrl("contract/" + contractAddress.Address)
	return httpapi.PostWithAPIResponseNoPayload[contract.DeployedContractMetaJSON](contractURL)
}
