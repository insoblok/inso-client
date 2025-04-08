package main

import (
	"context"
	ws "eth-toy-client/client/ws/wsplayground"
	"fmt"
	"github.com/ethereum/go-ethereum"
	accabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// Example of using GetAbiBin in a console application
func main() {
	logger := &ws.ConsoleLogger{}
	dir := "/Users/iyadi/playground/eth-toy-client/eth-toy-client/client/ws/wsmockusdc"
	contract := "MockUSDC"

	bin, parsedABI := ws.GetAbiBin(logger, dir, contract)
	logger.Logf("Binary: %x", bin)
	logger.Logf("ABI: %v", parsedABI)

	// 1. Connect to WS node
	client, err := ethclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatalf("❌ Failed to connect to WebSocket: %v", err)
	}
	logger.Logf("✅ Connected to Geth via WebSocket")

	// 3. Set up subscription
	logsCh := make(chan types.Log)
	query := ethereum.FilterQuery{
		// Add specific contract addresses if needed:
		// Addresses: []common.Address{common.HexToAddress("0xYourContract")},
	}

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to logs: %v", err)
	}

	logger.Logf("🎧 Listening for logs...")

	// 4. Loop and decode logs
	for {
		select {
		case err := <-sub.Err():
			log.Printf("⚠️ Subscription error: %v", err)

		case vLog := <-logsCh:
			handleLog(parsedABI, vLog)
		}
	}
}

func handleLog(parsedABI *accabi.ABI, vLog types.Log) {
	if len(vLog.Topics) == 0 {
		return
	}

	for name, event := range parsedABI.Events {
		if vLog.Topics[0] == event.ID {
			out := map[string]interface{}{}
			err := parsedABI.UnpackIntoMap(out, name, vLog.Data)
			if err != nil {
				log.Printf("❌ Failed to decode %s: %v", name, err)
				return
			}

			fmt.Printf("📢 Event: %s\n", name)
			fmt.Printf("   Block: %d\n", vLog.BlockNumber)
			fmt.Printf("   Tx: %s\n", vLog.TxHash.Hex())
			fmt.Printf("   Log: %s\n", vLog.Index)
			fmt.Printf("   Contract:%s\n", vLog.Address.Hex())
			for k, v := range out {
				fmt.Printf("   %s: %v\n", k, v)
			}
		}
	}
}
