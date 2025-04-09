package main

import (
	"context"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/devserver/devserver"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net/http"
	"time"
)

func main() {
	var port string
	var serverPort string
	flag.StringVar(&port, "port", "8545", "HTTP RPC port for the dev node")
	flag.StringVar(&serverPort, "serverPort", "8888", "HTTP RPC port for the supporting server")
	flag.Parse()

	devNodeConfig := devserver.DevNodeConfig{
		RPCPort: port,
	}
	rpcClient, ready, err := devserver.StartDevNode(devNodeConfig)
	if err != nil {
		log.Fatalf("Error starting dev node: %v", err)
	}

	select {
	case <-ready:
		log.Println("🚦 Node is ready. Proceed.")
	case <-time.After(5 * time.Second):
		log.Fatal("🕒 Timeout waiting for dev node to start.")
	}

	client := ethclient.NewClient(rpcClient)
	defer client.Close()

	var accounts []string
	err = rpcClient.Call(&accounts, "eth_accounts")
	if err != nil || len(accounts) == 0 {
		log.Fatalf("❌ Failed to get dev account: %v", err)
	}
	devAddr := common.HexToAddress(accounts[0])
	fmt.Printf("✅ Dev account: %s\n", devAddr.Hex())

	bal, err := client.BalanceAt(context.Background(), devAddr, nil)
	if err == nil {
		fmt.Printf("💰 Balance: %s wei\n", bal.String())
	}

	testAccount := devserver.LoadTestAccounts()
	fundedAccounts := devserver.FundTestAccounts(devAddr, rpcClient, testAccount)

	go func() {
		log.Println("🌐 Supporting HTTP server listening at http://localhost:" + serverPort + "...")
		err := http.ListenAndServe(
			":"+serverPort,
			devserver.SetupRoutes(contract.NewRegistry(), devAddr, port, fundedAccounts))
		if err != nil {
			log.Fatalf("❌ Failed to start HTTP server: %v", err)
		}
	}()

	log.Printf("📡 Dev node ready at http://localhost:%s — Press Ctrl+C to exit", port)
	select {}
}
