package contract

import (
	"context"
	"errors"
	"eth-toy-client/core/httpapi"
	toytypes "eth-toy-client/core/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"strings"
	"time"
)

// DeployContract sends a contract deployment tx via the dev server and waits for the contract address
func DeployContract(
	ctx context.Context,
	rpc *ethclient.Client,
	devServerURL string,
	fromAlias string,
	bytecode string,
) (common.Address, string, error) {
	// 🧹 Sanitize input
	data := strings.TrimPrefix(bytecode, "0x")
	data = strings.TrimSpace(data)
	if !strings.HasPrefix(data, "0x") {
		data = "0x" + data
	}

	// 📨 Build the request
	req := toytypes.SignTxRequest{
		From:  fromAlias,
		To:    "",
		Value: "0",
		Data:  data,
	}

	// 📤 POST to /api/send-tx
	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SendTxAPIResponse](
		devServerURL+"/api/send-tx", req,
	)
	if err != nil {
		return common.Address{}, "", fmt.Errorf("http error: %w", err)
	}
	if apiErr != nil {
		return common.Address{}, "", fmt.Errorf("api error: %s — %s", apiErr.Code, apiErr.Message)
	}

	txHash := common.HexToHash(apiResp.TxHash)
	fmt.Printf("🚀 Deployment tx sent: %s\n", txHash.Hex())

	// ⏳ Wait for receipt
	for i := 0; i < 60; i++ {
		receipt, err := rpc.TransactionReceipt(ctx, txHash)
		if err == nil && receipt != nil {
			if receipt.ContractAddress == (common.Address{}) {
				return common.Address{}, txHash.Hex(), errors.New("receipt has no contract address")
			}
			fmt.Printf("✅ Contract deployed at %s (block %d)\n",
				receipt.ContractAddress.Hex(), receipt.BlockNumber.Uint64())
			return receipt.ContractAddress, txHash.Hex(), nil
		}
		time.Sleep(1 * time.Second)
	}

	return common.Address{}, txHash.Hex(), fmt.Errorf("timeout waiting for contract deployment")
}
