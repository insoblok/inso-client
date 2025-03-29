package main

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func main() {
	// Connect to the running dev node via HTTP
	client, err := ethclient.Dial("http://localhost:8565")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum node:", err)
	}
	defer client.Close()

	num, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("Failed to connect to Ethereum node:", err)
	}

	log.Printf("Max Block available number: %d", num+1)

	for i := uint64(0); i <= num; i++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			log.Fatalf("Failed to fetch block %d: %v", i, err)
		}
		log.Printf(
			"Block %d: Hash: %s Time: %s\n Transactions: %d",
			i,
			block.Hash().Hex(),
			time.Unix(int64(block.Time()), 0).Format(time.RFC3339),
			block.Transactions().Len())
	}
}
