package devserver

import (
	"context"
	"eth-toy-client/config"
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

var ServerName = config.Servers.DevServer

func TestDeployContract(t *testing.T) {
	serverConfig := ServerName.GetServerConfig()
	apiSendTxURL := serverConfig.GetServerUrl("api/deploy-contract")

	req := toytypes.SignTxRequest{
		From:  "alice",
		To:    "",
		Value: "0",
		Data:  devserver.MockusdcMetaData.Bin[2:],
	}

	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SendTxAPIResponse](
		apiSendTxURL,
		req,
	)
	require.NoError(t, err, "❌ Error sending tx")
	require.Nil(t, apiErr, "❌ Error sending tx")
	require.NotNil(t, apiResp, "❌ Error sending tx")
	require.NotEmpty(t, apiResp.TxHash, "❌ TxHash is empty")

	t.Logf("✅ℹ️Received Tx Hash: %s\n", apiResp.TxHash)

	client, err := serverConfig.DevNodeConfig.GetEthClient()
	require.NoError(t, err, "❌ Error connecting to dev node")
	ctx := context.Background()

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

	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	require.NoError(t, err, "❌ failed to fetch contract code")
	require.NotNil(t, code, "❌ code is nil")
	logutil.Infof("ℹ️contract code: %x", string(code))
	require.True(t, len(code) > 0, "❌ empty contract code")

}
