package devserver

import (
	"context"
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	devserver "eth-toy-client/servers/devserver/devserver/test/contracts/mockusdc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func GetNonce(serverConfig config.ServerConfig, alias string) (*toytypes.PendingNonceResponse, *httpapi.APIError, error) {
	pendingNonceUrl := serverConfig.GetServerUrl("api/pending-nonce")
	req := toytypes.PendingNonceRequest{
		Alias: alias,
	}
	return httpapi.PostWithAPIResponse[toytypes.PendingNonceResponse](
		pendingNonceUrl,
		req,
	)
}

func RegisterContractAddress(alias string, serverConfig config.ServerConfig, address string, abi string) (*toytypes.AliasRegisterResponse, *httpapi.APIError, error) {
	registerContractURL := serverConfig.GetServerUrl("api/register-contract")
	req := contract.DeployedContractMetaJSON{
		Alias:   alias,
		Address: address,
		ABI:     abi,
	}
	return httpapi.PostWithAPIResponse[toytypes.AliasRegisterResponse](
		registerContractURL,
		req,
	)
}

func TestDeployContract(t *testing.T) {
	devServerConfig := config.Servers.DevServer.GetServerConfig()
	alias := "alice"

	pendingNonce, pendingApiError, err := GetNonce(devServerConfig, alias)
	require.NoError(t, err, "❌ Failed to get pending nonce")
	require.Nil(t, pendingApiError, "❌ APIError Failed to get pending nonce")
	require.NotNil(t, pendingNonce, "❌ PendingNonce is nil")

	logServerConfig := config.Servers.LogServer.GetServerConfig()
	aliasRegisterResponse, aliasApiError, err := RegisterContractAddress(alias, logServerConfig, pendingNonce.Address, devserver.MockusdcMetaData.ABI)
	require.NoError(t, err, "❌ Failed to register contract address")
	require.Nil(t, aliasApiError, "❌ APIError Failed to register contract address")
	require.NotNil(t, aliasRegisterResponse, "❌ AliasRegisterResponse is nil")

	apiSendTxURL := devServerConfig.GetServerUrl("api/deploy-contract")

	req := toytypes.DeployContractRequest{
		From:  alias,
		Nonce: pendingNonce.Nonce,
		Data:  devserver.MockusdcMetaData.Bin[2:],
	}

	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.ContractDeploymentResponse](
		apiSendTxURL,
		req,
	)

	require.NoError(t, err, "❌ Error sending tx")
	require.Nil(t, apiErr, "❌ Error sending tx")
	require.NotNil(t, apiResp, "❌ Error sending tx")
	require.NotEmpty(t, apiResp.TxHash, "❌ TxHash is empty")
	require.NotEmpty(t, apiResp.ExpectedContractAddress)

	t.Logf("✅ℹ️Received Tx Hash: %s\n", apiResp.TxHash)
	t.Logf("✅ℹ️Expected Contract Address: %s\n", apiResp.ExpectedContractAddress)

	require.NoError(t, err, "❌ Error connecting to dev node")
	ctx := context.Background()
	client, err := devServerConfig.DevNodeConfig.GetEthClient()
	var receipt *types.Receipt
	for i := 0; i < 60; i++ {
		receipt, err = client.TransactionReceipt(ctx, common.HexToHash(apiResp.TxHash))
		if err == nil && receipt != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	require.NotNil(t, receipt, "❌🧾 ⏱️ timeout waiting for tx %s", apiResp.TxHash)

	t.Log("ℹ️Transaction Receipt Details:")
	t.Logf("  ℹ️Status: %d\n", receipt.Status) // Status: 1 (success) or 0 (failure)
	t.Logf("  ℹ️Transaction Hash: %s\n", receipt.TxHash.Hex())
	t.Logf("  ℹ️Contract Address: %s\n", receipt.ContractAddress.Hex())
	t.Logf("  ℹ️Block Number: %d\n", receipt.BlockNumber.Uint64())
	t.Logf("  ℹ️Gas Used: %d\n", receipt.GasUsed)
	t.Logf("  ℹ️Logs:")
	for i, log := range receipt.Logs {
		t.Logf("    ℹ️Log #%d: %+v\n", i, log)
	}

	require.Equal(t, uint64(1), receipt.Status, "❌ transaction failed, status: %d", receipt.Status)
	require.Equal(t, apiResp.ExpectedContractAddress, receipt.ContractAddress.Hex(), "❌ contract address mismatch")

	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	require.NoError(t, err, "❌ failed to fetch contract code")
	require.NotNil(t, code, "❌ code is nil")
	logutil.Infof("ℹ️contract code: %x", string(code))
	require.True(t, len(code) > 0, "❌ empty contract code")

}
