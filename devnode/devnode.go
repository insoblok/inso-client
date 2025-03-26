package main

import (
	"context"
	"flag"
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
	// 🎛️ CLI flag for port
	var port string
	flag.StringVar(&port, "port", "8545", "HTTP RPC port for the dev node")
	flag.Parse()

	var gethCmd = "/Users/iyadi/github/ethereum/go-ethereum/build/bin/geth" // leave this hardcoded for now

	// 🚀 Start geth
	cmd := exec.Command(gethCmd,
		"--dev",
		"--http",
		"--http.api", "eth,net,web3,personal",
		"--http.addr", "127.0.0.1",
		"--http.port", port,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("🚀 Starting Geth dev node on port %s...", port)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("❌ Failed to start geth: %v", err)
	}
	defer func() {
		log.Println("🛑 Shutting down Geth...")
		cmd.Process.Kill()
	}()

	// ⏳ Wait for readiness
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

	// 🧙 Query dev account
	var accounts []string
	err = rpcClient.Call(&accounts, "eth_accounts")
	if err != nil || len(accounts) == 0 {
		log.Fatalf("❌ Failed to get dev account: %v", err)
	}
	devAddr := common.HexToAddress(accounts[0])
	fmt.Printf("✅ Dev account: %s\n", devAddr.Hex())

	// 💰 Query balance
	bal, err := client.BalanceAt(context.Background(), devAddr, nil)
	if err == nil {
		fmt.Printf("💰 Balance: %s wei\n", bal.String())
	}

	log.Printf("📡 Dev node ready at http://localhost:%s — Press Ctrl+C to exit", port)
	select {}
}
