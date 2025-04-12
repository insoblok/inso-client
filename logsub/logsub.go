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
	Decoder     Decoder
}

func NewLogListener(config *ListenerConfig, broadcaster logbus.LogBroadcaster, decoder Decoder) (*LogListener, error) {
	client, err := ethclient.Dial(config.WebSocketURL)
	if err != nil {
		return nil, err
	}

	return &LogListener{
		Config:      config,
		Broadcaster: broadcaster,
		Client:      client,
		Decoder:     decoder,
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

			// Decode the log using the DefaultDecoder
			event, err := l.Decoder.DecodeLog(logEvent)
			if err != nil {
				log.Printf("❌ Failed to decode log: %v", err)
				continue
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

func (l *LogListener) Listen2(ctx context.Context) {
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

// Decoder interface that needs to be implemented for decoding logs
type Decoder interface {
	DecodeLog(log types.Log) (logbus.LogEvent, error) // Decodes a log into a LogEvent
}

// DefaultDecoder is the default implementation of the Decoder interface
type DefaultDecoder struct{}

// DecodeLog is the default method that simply converts the log into a generic LogEvent
func (d *DefaultDecoder) DecodeLog(log types.Log) (logbus.LogEvent, error) {
	// For now, assume generic event with no ABI decoding
	event := logbus.LogEvent{
		Contract:  log.Address.Hex(),
		Event:     "UnknownEventForNow", // Assume an unknown event
		TxHash:    log.TxHash.Hex(),
		Block:     log.BlockNumber,
		Timestamp: time.Now().Unix(),
		Args:      make(map[string]interface{}), // Empty args for now
	}

	// Return the generic event
	return event, nil
}
