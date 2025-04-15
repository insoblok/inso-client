package logserver

import (
	"context"
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/servers/servers"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"net/http"
)

type LogServer struct{}

func (logServer *LogServer) Name() config.ServerName {
	return "LogServer"
}

func (logServer *LogServer) InitService(nodeClient *servers.NodeClient, serverConfig config.ServerConfig) (config.ServerConfig, http.Handler) {
	contractRegistry := contract.NewRegistry()
	go InitLogListener(nodeClient)

	handlers := SetupRoutes(serverConfig, contractRegistry)
	return serverConfig, handlers
}

func InitLogListener(nodeClient *servers.NodeClient) {
	logsCh := make(chan types.Log)
	query := ethereum.FilterQuery{
		// Add specific contract addresses if needed:
		// Addresses: []common.Address{common.HexToAddress("0xYourContract")},
	}

	sub, err := nodeClient.WSClient.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to logs: %v", err)
	}
	log.Println("🎧 Listening for logs...")
	for {
		select {
		case err := <-sub.Err():
			log.Printf("⚠️ Subscription error: %v", err)

		case vLog := <-logsCh:
			log.Printf("📄 Received log: %+v", vLog)
		}
	}
}
