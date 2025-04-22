package logserver

import (
	"context"
	"eth-toy-client/config"
	contract "eth-toy-client/core/contracts"
	toytypes "eth-toy-client/core/types"
	"eth-toy-client/logbus"
	"eth-toy-client/logsub"
	"eth-toy-client/servers/servers"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"net/http"
)

type LogServer struct{}

type LogDecoder struct {
	registry *contract.Registry
}

func (logDecoder *LogDecoder) DecodeLog(logEvent types.Log) (logbus.LogEvent, error) {
	fmt.Printf("Decodong log: %s\n", logEvent)
	contractAddr := toytypes.ContractAddress{Address: logEvent.Address.Hex()}
	evt := logsub.DecodeGenericLog(logEvent)
	info, ok := logDecoder.registry.Get(contractAddr)
	if !ok {
		log.Printf("No info about contract %s", contractAddr)
		return evt, nil
	}
	if len(logEvent.Topics) == 0 {
		log.Printf("No topics in log %s", logEvent)
		return evt, nil
	}

	fmt.Printf("TopicsLen: %v\n", len(logEvent.Topics))

	for name, event := range info.ParsedABI.Events {
		if logEvent.Topics[0] == event.ID {
			out := map[string]interface{}{}
			err := info.ParsedABI.UnpackIntoMap(out, name, logEvent.Data)
			if err != nil {
				log.Printf("❌ Failed to decode %s: %v", name, err)

			} else {
				fmt.Printf("📢 Event: %s\n", name)
				fmt.Printf("   Block: %d\n", logEvent.BlockNumber)
				fmt.Printf("   BlockHash: %v\n", logEvent.BlockHash.Hex())
				fmt.Printf("   Tx: %s\n", logEvent.TxHash.Hex())
				fmt.Printf("   LogIndex: %d\n", logEvent.Index)
				fmt.Printf("   Contract:%s\n", logEvent.Address.Hex())

				fmt.Printf("🔍 Event Structure Details for '%s':\n", name)
				eventDetails := event.Inputs

				indexedArgs := make([]abi.Argument, 0)
				noneIndexedArgs := make([]abi.Argument, 0)

				for i, input := range eventDetails {
					fmt.Printf("%d   Name: %s, Type: %s, Indexed: %v\n", i, input.Name, input.Type.String(), input.Indexed)
					if input.Indexed {
						indexedArgs = append(indexedArgs, input)
					} else {
						noneIndexedArgs = append(noneIndexedArgs, input)
					}
				}

				if len(noneIndexedArgs) != len(out) {
					log.Printf("⚠️ Mismatch between noneIndexedArgs length (%d) and out length (%d)", len(noneIndexedArgs), len(out))
				}

				if len(indexedArgs) != (len(logEvent.Topics) - 1) {
					log.Printf("⚠️ Mismatch between IdexedArgs length (%d) and topics (%d)", len(indexedArgs), len(logEvent.Topics))
				}

				indexedValues := logEvent.Topics[1:]
				for i, input := range indexedArgs {
					fmt.Printf("  %s : %s\n", input.Name, indexedValues[i].Hex())
				}

				for k, v := range out {
					fmt.Printf("   %s: %v\n", k, v)
				}
			}
		}
	}

	return logbus.LogEvent{}, nil
}

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
	logDecoder := &LogDecoder{
		registry: contractRegistry,
	}

	go InitLogListener(nodeClient, broadcaster, logDecoder)

	handlers := SetupRoutes(serverConfig, contractRegistry)
	return serverConfig, handlers
}

func InitLogListener(nodeClient *servers.NodeClient, broadcaster logbus.LogBroadcaster, decoder logsub.Decoder) {
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
