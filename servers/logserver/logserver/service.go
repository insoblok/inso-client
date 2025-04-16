package logserver

import (
	"context"
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	"eth-toy-client/logbus"
	"eth-toy-client/logsub"
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
	broadcaster := logbus.NewLogBroadcaster()
	eventBus := make(chan logbus.LogEvent, 10)

	consoleConsumer := &ConsoleConsumer{
		Name:             "ConsoleConsumer",
		ContractRegistry: contractRegistry,
		Events:           eventBus}
	go consoleConsumer.Consume()

	broadcaster.Subscribe(eventBus)
	decoder := &logsub.DefaultDecoder{
		DecoderFn: logsub.GetDecoderFn(),
	}
	go InitLogListener(nodeClient, broadcaster, decoder)

	handlers := SetupRoutes(serverConfig, contractRegistry)
	return serverConfig, handlers
}

func InitLogListener(nodeClient *servers.NodeClient, broadcaster logbus.LogBroadcaster, decoder *logsub.DefaultDecoder) {
	logsCh := make(chan types.Log)
	query := ethereum.FilterQuery{}

	sub, err := nodeClient.WSClient.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to logs: %v", err)
	}
	log.Println("🎧 Listening for logs...")
	for {
		select {
		case err := <-sub.Err():
			log.Printf("⚠️ Subscription error: %v", err)

		case logEvent := <-logsCh:
			//log.Printf("📄 Received log: %+v", logEvent)
			event, errX := decoder.DecodeLog(logEvent)
			if errX != nil {
				log.Printf("❌ Failed to decode log: %v", err)
				continue
			}
			broadcaster.Publish(event)

		}
	}
}
