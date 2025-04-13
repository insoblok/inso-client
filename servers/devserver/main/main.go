package main

import (
	"context"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/servers/devserver/devserver"
	"eth-toy-client/servers/servers"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"net/http"
)

func main() {
	serverConfig, nodeClient := servers.EstablishConnectionToDevNode()
	defer nodeClient.Client.Close()
	defer nodeClient.RPCClient.Close()

	var accounts []string
	err := nodeClient.RPCClient.Call(&accounts, "eth_accounts")
	if err != nil || len(accounts) == 0 {
		log.Fatalf("❌ Failed to get dev account: %v", err)
	}
	devAddr := common.HexToAddress(accounts[0])
	fmt.Printf("✅ Dev account: %s\n", devAddr.Hex())

	bal, err := nodeClient.Client.BalanceAt(context.Background(), devAddr, nil)
	if err == nil {
		fmt.Printf("💰 Balance: %s wei\n", bal.String())
	}

	testAccount := devserver.LoadTestAccounts()
	fundedAccounts := devserver.FundTestAccounts(devAddr, nodeClient.RPCClient, testAccount)

	go func() {
		log.Println("🌐 " + serverConfig.Name + "listening at http://localhost:" + serverConfig.Port + "...")
		err := http.ListenAndServe(
			":"+serverConfig.Port,
			devserver.SetupRoutes(contract.NewRegistry(), devAddr, nodeClient, fundedAccounts))
		if err != nil {
			log.Fatalf("❌ Failed to start HTTP server: %v", err)
		}
	}()

	select {}
}
