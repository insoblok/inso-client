package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// Connect to the running dev node via HTTP
	client, err := ethclient.Dial("http://localhost:8565")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum node:", err)
	}
	defer client.Close()

	// Fetch latest block number
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Failed to get block number:", err)
	}

	fmt.Printf("🎉 Connected to Ethereum dev node! Current block number: %d\n", blockNumber)
}
