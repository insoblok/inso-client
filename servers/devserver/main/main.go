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
	devServer := &DevServer{}
	servers.StartMicroService(devServer)
	select {}
}

type DevServer struct{}

func (devServer *DevServer) Name() string {
	return "DevServer"
}

func (devServer *DevServer) InitService(nodeClient *servers.NodeClient, serverConfig servers.ServerConfig) (servers.ServerConfig, http.Handler) {
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
	} else {
		log.Fatalf("❌ Failed to obtain balance for devAddr: %v", err)
	}

	testAccount := devserver.LoadTestAccounts()
	fundedAccounts := devserver.FundTestAccounts(devAddr, nodeClient.RPCClient, testAccount)
	handler := devserver.SetupRoutes(contract.NewRegistry(), devAddr, nodeClient, fundedAccounts)
	return serverConfig, handler
}
