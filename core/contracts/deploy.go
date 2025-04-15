package contract

import (
	"context"
	"eth-toy-client/core/httpapi"
	"eth-toy-client/core/logutil"
	toytypes "eth-toy-client/core/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
	"time"
)

func DeployContract(
	ctx context.Context,
	client *ethclient.Client,
	serverURL string,
	fromAlias string,
	bytecode string,
) (common.Address, string, error) {
	req := toytypes.SignTxRequest{
		From:  fromAlias,
		To:    "",
		Value: "0",
		Data:  bytecode,
	}

	apiResp, apiErr, err := httpapi.PostWithAPIResponse[toytypes.SendTxAPIResponse](
		serverURL+"/api/send-tx",
		req,
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

	// Print receipt details for debugging
	fmt.Println("Transaction Receipt Details:")
	fmt.Printf("  Status: %d\n", receipt.Status) // Status: 1 (success) or 0 (failure)
	fmt.Printf("  Transaction Hash: %s\n", receipt.TxHash.Hex())
	fmt.Printf("  Contract Address: %s\n", receipt.ContractAddress.Hex())
	fmt.Printf("  Block Number: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
	fmt.Println("  Logs:")
	for i, log := range receipt.Logs {
		fmt.Printf("    Log #%d: %+v\n", i, log)
	}

	if receipt.Status != 1 {
		return common.Address{}, txHash, fmt.Errorf("transaction failed, status: %d", receipt.Status)
	}

	logutil.Infof("contract address: %s", receipt.ContractAddress.Hex())

	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	if err != nil {
		return common.Address{}, txHash, fmt.Errorf("failed to fetch contract code: %w", err)
	}
	logutil.Infof("contract code: %x", string(code))
	logutil.Infof("contract Tx Hash: %s", receipt.TxHash)
	if len(code) == 0 {
		return common.Address{}, txHash, fmt.Errorf("contract code is empty — deployment likely failed")
	}

	return receipt.ContractAddress, txHash, nil
}

type AliasDeployRequest struct {
	From     string `json:"from"`          // e.g., "alice"
	Alias    string `json:"alias"`         // e.g., "counter-v1"
	Bytecode string `json:"bytecode"`      // "0x..."
	ABI      string `json:"abi,omitempty"` // optional for now
}

type ContractMeta struct {
	Alias     string         `json:"alias"`
	Address   common.Address `json:"address"`
	TxHash    common.Hash    `json:"txHash"`
	ABI       string         `json:"abi"`
	Timestamp time.Time      `json:"timestamp"`
}

type Registry struct {
	mu      sync.RWMutex
	entries map[toytypes.ContractAddress]DeployedContractMetaJSON
}

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[toytypes.ContractAddress]DeployedContractMetaJSON),
	}
}

func (r *Registry) Add(meta DeployedContractMetaJSON) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	contractAddress := toytypes.ContractAddress{Address: meta.Address}
	if _, exists := r.entries[contractAddress]; exists {
		return fmt.Errorf("ContractAddress already exists: %s", contractAddress.Address)
	}

	r.entries[contractAddress] = meta
	return nil
}

func (r *Registry) All() []DeployedContractMetaJSON {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make([]DeployedContractMetaJSON, 0, len(r.entries))
	for _, meta := range r.entries {
		entries = append(entries, meta)
	}
	return entries
}

func (r *Registry) Get(contractAddress toytypes.ContractAddress) (DeployedContractMetaJSON, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	meta, ok := r.entries[contractAddress]
	return meta, ok
}

type AliasDeployResponse struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
	TxHash  string `json:"txHash"`
}

type DeployedContractMetaJSON struct {
	Alias     string `json:"alias"`
	Address   string `json:"address"`
	TxHash    string `json:"txHash"`
	ABI       string `json:"abi"`
	Bytecode  string `json:"bytecode"`
	Timestamp int64  `json:"timestamp"`
	Owner     string `json:"owner"`
	Overwrite bool   `json:"overwrite"`
}
