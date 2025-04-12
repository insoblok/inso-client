package logsub

import (
	"context"
	"eth-toy-client/logbus"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"
)

type ListenerConfig struct {
	WebSocketURL string // WebSocket URL to connect to the Ethereum node
}

type LogListener struct {
	Config      *ListenerConfig
	Broadcaster logbus.LogBroadcaster
	Client      *ethclient.Client
}

func NewLogListener(config *ListenerConfig, broadcaster logbus.LogBroadcaster) (*LogListener, error) {
	client, err := ethclient.Dial(config.WebSocketURL)
	if err != nil {
		return nil, err
	}

	return &LogListener{
		Config:      config,
		Broadcaster: broadcaster,
		Client:      client,
	}, nil
}

func (l *LogListener) Listen(ctx context.Context) {
	log.Printf("✅ LogListener started. Listening for logs...")

	// Example log filter: change as needed
	filter := ethereum.FilterQuery{
		Addresses: []common.Address{}, // Empty list means all contracts
	}

	logs := make(chan types.Log)

	sub, err := l.Client.SubscribeFilterLogs(ctx, filter, logs)
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to logs: %v", err)
	}

	for {
		select {
		case logEvent := <-logs:
			// Print the received log to console
			log.Printf("🎤 Received Log: %v", logEvent)

			// Convert the Ethereum log into a LogEvent
			event := logbus.LogEvent{
				Contract:  logEvent.Address.Hex(),
				Event:     "UnknownEventForNow", // For now, assume generic event
				TxHash:    logEvent.TxHash.Hex(),
				Block:     logEvent.BlockNumber,
				Timestamp: time.Now().Unix(),
				Args:      make(map[string]interface{}), // Just an empty map for now
			}

			// Publish this log event to LogBroadcaster
			l.Broadcaster.Publish(event)
		case err := <-sub.Err():
			log.Printf("⚠️ Subscription error: %v", err)

		case <-ctx.Done():
			log.Println("👋 LogListener exiting...")
			return
		}
	}
}

func (l *LogListener) StartSimulateListening() {
	// Here you would start the actual logic for listening to the Ethereum logs
	// This will involve setting up log filters, subscribing to logs, etc.

	log.Println("✅LogListener started. Listening for logs...")

	event := logbus.LogEvent{
		Event:  "TestEvent",
		TxHash: "0x123",
	}

	l.Broadcaster.Publish(event)

}

type PrintToConsole struct {
	Name   string
	Events chan logbus.LogEvent
}

func (p *PrintToConsole) Consume() {
	for event := range p.Events {
		log.Printf("🚀 %s received event: %s with TxHash: %s", p.Name, event.Event, event.TxHash)
	}
}
