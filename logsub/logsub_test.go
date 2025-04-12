package logsub

import (
	"eth-toy-client/logbus"
	"log"
	"testing"
	"time"
)

func TestLogListenerWithPrintToConsole(t *testing.T) {
	broadcaster := logbus.NewLogBroadcaster()

	config := &ListenerConfig{
		WebSocketURL: "ws://localhost:8546", // Change to actual WebSocket URL for the DevNode
	}

	listener, err := NewLogListener(config, broadcaster)
	if err != nil {
		t.Fatalf("Error creating LogListener: %v", err)
	}

	consoleChannel := make(chan logbus.LogEvent, 10)
	consoleConsumer := &PrintToConsole{Name: "ConsoleConsumer", Events: consoleChannel}

	go consoleConsumer.Consume()

	broadcaster.Subscribe(consoleChannel)

	go listener.StartListening()

	time.Sleep(5 * time.Second)

	log.Printf("🐳 Done")
}
