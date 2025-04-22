package contract

import (
	"context"
	"eth-toy-client/core/httpapi"
	toytypes "eth-toy-client/core/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

func DeployContract(
	ctx context.Context,
	client *ethclient.Client,
	serverURL string,
	fromAlias string,
	bytecode string,
) (common.Address, string, error) {
	// 📨 Compose request
	req := toytypes.SignTxRequest{
		From:  fromAlias,
		To:    "",
		Value: "0",
		Data:  bytecode,
	}

	// 📤 Send tx
	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SendTxAPIResponse](
		serverURL+"/api/send-tx", req,
	)
	if err != nil {
		return common.Address{}, "", fmt.Errorf("http error: %w", err)
	}
	if apiErr != nil {
		return common.Address{}, "", fmt.Errorf("api error: %s — %s", apiErr.Code, apiErr.Message)
	}
	if apiResp == nil || apiResp.TxHash == "" {
		return common.Address{}, "", fmt.Errorf("unexpected: nil or empty tx response")
	}

	txHash := apiResp.TxHash
	fmt.Printf("🚀 Deployment tx sent: %s\n", txHash)

	// ⏳ Wait for receipt
	var receipt *types.Receipt
	for i := 0; i < 60; i++ {
		receipt, err = client.TransactionReceipt(ctx, common.HexToHash(txHash))
		if err == nil && receipt != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if receipt == nil {
		return common.Address{}, txHash, fmt.Errorf("⏱️ timeout waiting for tx %s", txHash)
	}

	// 🧪 Check deployed contract code
	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	if err != nil {
		return common.Address{}, txHash, fmt.Errorf("failed to fetch contract code: %w", err)
	}
	if len(code) == 0 {
		return common.Address{}, txHash, fmt.Errorf("contract code is empty — deployment likely failed")
	}

	return receipt.ContractAddress, txHash, nil
}
