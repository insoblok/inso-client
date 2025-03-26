package main

import (
	"context"
	"encoding/json"
	"eth-toy-client/devnode"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	var port string
	var serverPort string
	flag.StringVar(&port, "port", "8545", "HTTP RPC port for the dev node")
	flag.StringVar(&serverPort, "serverPort", "8888", "HTTP RPC port for the supporting server")
	flag.Parse()

	var gethCmd = "/Users/iyadi/github/ethereum/go-ethereum/build/bin/geth"

	cmd := exec.Command(gethCmd,
		"--dev",
		"--http",
		"--http.api", "eth,net,web3,txpool,miner,admin,debug",
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

	devnode.LoadTestAccounts()
	devnode.FundTestAccounts(devAddr, rpcClient)

	// ✅ ✅ ✅ START HTTP SERVER
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/dev-account", func(w http.ResponseWriter, r *http.Request) {
			resp := struct {
				Address string `json:"address"`
			}{
				Address: devAddr.Hex(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		})

		log.Println("🌐 Supporting HTTP server listening at http://localhost:" + serverPort + "...")
		err := http.ListenAndServe(":"+serverPort, mux)
		if err != nil {
			log.Fatalf("❌ Failed to start HTTP server: %v", err)
		}
	}()
	// ✅ ✅ ✅ END HTTP SERVER

	log.Printf("📡 Dev node ready at http://localhost:%s — Press Ctrl+C to exit", port)
	select {}
}
