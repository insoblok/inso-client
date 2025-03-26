package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	// 1. Start geth as a child process
	var port = "8565"
	var gethCmd = "/Users/iyadi/github/ethereum/go-ethereum/build/bin/geth"

	cmd := exec.Command(gethCmd,
		"--dev",
		"--http",
		"--http.api", "eth,net,web3,personal",
		"--http.addr", "127.0.0.1",
		"--http.port", port,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("🚀 Starting Geth dev node...")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("❌ Failed to start geth: %v", err)
	}

	defer func() {
		log.Println("🛑 Shutting down Geth...")
		cmd.Process.Kill()
	}()

	// 2. Wait for geth to be ready (simple poll)
	var rpcClient *rpc.Client
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		rpcClient, err = rpc.Dial("http://localhost:" + port)
		if err == nil {
			break
		}
		log.Println("⏳ Waiting for Geth to be ready...")
	}
	if rpcClient == nil {
		log.Fatal("❌ Geth did not start in time")
	}
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)
	defer client.Close()

	// 3. Query dev account
	var accounts []string
	err = rpcClient.Call(&accounts, "eth_accounts")
	if err != nil || len(accounts) == 0 {
		log.Fatalf("❌ Failed to get dev account: %v", err)
	}
	devAddr := common.HexToAddress(accounts[0])
	fmt.Printf("✅ Dev account: %s\n", devAddr.Hex())

	// 4. (Optional) Get balance
	bal, err := client.BalanceAt(context.Background(), devAddr, nil)
	if err == nil {
		fmt.Printf("💰 Balance: %s wei\n", bal.String())
	}

	log.Println("📡 Dev node is ready. Press Ctrl+C to exit.")
	select {} // Block forever
}
