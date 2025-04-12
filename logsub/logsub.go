package logsub

import (
	"eth-toy-client/logbus"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
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

func (l *LogListener) StartListening() {
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
