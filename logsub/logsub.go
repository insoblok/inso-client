package logsub

import (
	"context"
	"eth-toy-client/logbus"
	"fmt"
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

type Decoder interface {
	DecodeLog(log types.Log) (logbus.LogEvent, error) // Decodes a log into a LogEvent
}

type DefaultDecoder struct {
	DecoderFn func(logType logbus.LogType) (Decoder, error)
}

func (d *DefaultDecoder) DecodeLog(log types.Log) (logbus.LogEvent, error) {
	logType := logbus.GetLogType(log)

	decoder, err := d.DecoderFn(logType)
	if err != nil {
		return d.decodeGenericLog(log)
	}

	return decoder.DecodeLog(log)
}

func (d *DefaultDecoder) decodeGenericLog(log types.Log) (logbus.LogEvent, error) {
	return logbus.LogEvent{
		Contract:  log.Address.Hex(),
		Event:     "UnknownEvent",
		TxHash:    log.TxHash.Hex(),
		Block:     log.BlockNumber,
		Timestamp: time.Now().Unix(),
		Args:      make(map[string]interface{}),
	}, nil
}

func NewDefaultDecoder(decoderFn func(logType logbus.LogType) (Decoder, error)) *DefaultDecoder {
	return &DefaultDecoder{
		DecoderFn: decoderFn,
	}
}

type DecoderRegistry struct {
	Decoders map[logbus.LogType]Decoder
}

func NewDecoderRegistry() *DecoderRegistry {
	return &DecoderRegistry{
		Decoders: make(map[logbus.LogType]Decoder),
	}
}

func (r *DecoderRegistry) RegisterDecoder(logType logbus.LogType, decoder Decoder) {
	r.Decoders[logType] = decoder
}

func (r *DecoderRegistry) GetDecoderFn() func(logType logbus.LogType) (Decoder, error) {
	return func(logType logbus.LogType) (Decoder, error) {
		decoder, exists := r.Decoders[logType]
		if !exists {
			return &DefaultDecoder{}, fmt.Errorf("no decoder found for log type: %v", logType)
		}
		return decoder, nil
	}
}

func GetDecoderFn() func(logType logbus.LogType) (Decoder, error) {
	decoderRegistry := NewDecoderRegistry()
	decoderRegistry.RegisterDecoder(logbus.TransactionLog, &TransactionLogDecoder{})
	decoderRegistry.RegisterDecoder(logbus.UnknownEventLog, &UnknownEventLogDecoder{})
	return decoderRegistry.GetDecoderFn()
}

type TransactionLogDecoder struct{}
type EventLogDecoder struct{}
type TokenTransferLogDecoder struct{}
type CustomContractEventDecoder struct{}
type StateChangeLogDecoder struct{}
type ErrorLogDecoder struct{}
type InternalTransactionDecoder struct{}
type UnknownEventLogDecoder struct{}

func (t *TransactionLogDecoder) DecodeLog(log types.Log) (logbus.LogEvent, error) {
	return logbus.LogEvent{
		LogType:   logbus.TransactionLog,
		Contract:  log.Address.Hex(),
		Event:     "TransactionLog",
		TxHash:    log.TxHash.Hex(),
		Block:     log.BlockNumber,
		Timestamp: time.Now().Unix(),
		Args:      make(map[string]interface{}),
	}, nil
}

func (t *UnknownEventLogDecoder) DecodeLog(log types.Log) (logbus.LogEvent, error) {
	return logbus.LogEvent{
		LogType:   logbus.UnknownEventLog,
		Contract:  log.Address.Hex(),
		Event:     "UnknownEventLog",
		TxHash:    log.TxHash.Hex(),
		Block:     log.BlockNumber,
		Timestamp: time.Now().Unix(),
		Args:      make(map[string]interface{}),
	}, nil
}
