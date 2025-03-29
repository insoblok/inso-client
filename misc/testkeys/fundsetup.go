package testkeys

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"time"
)

func main() {
	ctx := context.Background()

	// Connect to the dev node
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("❌ Failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	rpcClient, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("❌ Failed to connect via RPC: %v", err)
	}
	defer rpcClient.Close()

	// Fetch dev account from eth_accounts
	var devAccounts []string
	err = rpcClient.Call(&devAccounts, "eth_accounts")
	if err != nil || len(devAccounts) == 0 {
		log.Fatalf("❌ Failed to fetch dev account: %v", err)
	}
	devAddr := common.HexToAddress(devAccounts[0])
	log.Printf("🧙 Using dev account: %s", devAddr.Hex())

	// Get chain ID (usually 1337 in dev mode)
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to get chain ID: %v", err)
	}

	// Loop through test accounts and fund each
	for _, acct := range TestAccounts {
		log.Printf("💸 Funding %s (%s)...", acct.Name, acct.Addr.Hex())

		// Get dev's latest nonce
		nonce, err := client.PendingNonceAt(ctx, devAddr)
		if err != nil {
			log.Fatalf("❌ Failed to get nonce: %v", err)
		}

		// Suggest gas parameters
		tip, _ := client.SuggestGasTipCap(ctx)
		feeCap := new(big.Int).Add(tip, big.NewInt(2e9))

		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: tip,
			GasFeeCap: feeCap,
			Gas:       21000,
			To:        &acct.Addr,
			Value:     big.NewInt(1e18), // 1 ETH
		})

		// Because dev account is unlocked, we can send directly
		err = client.SendTransaction(ctx, tx)
		if err != nil {
			log.Fatalf("❌ Failed to send tx: %v", err)
		}

		log.Printf("✅ Sent funding tx: %s", tx.Hash().Hex())
	}

	log.Println("⏳ Waiting for blocks to confirm...")
	time.Sleep(2 * time.Second)

	// Optional: confirm balances
	for _, acct := range TestAccounts {
		bal, err := client.BalanceAt(ctx, acct.Addr, nil)
		if err != nil {
			log.Printf("⚠️  Failed to get balance: %v", err)
			continue
		}
		log.Printf("💰 %s: %s wei", acct.Name, bal.String())
	}
}
