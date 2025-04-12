package logsub

import (
	"context"
	"eth-toy-client/logbus"
	"log"
	"testing"
	"time"
)

func TestLogSimulateListenerWithPrintToConsole(t *testing.T) {
	broadcaster := logbus.NewLogBroadcaster()

	config := &ListenerConfig{
		WebSocketURL: "ws://localhost:8546", // Change to actual WebSocket URL for the DevNode
	}

	listener, err := NewLogListener(config, broadcaster, &DefaultDecoder{})
	if err != nil {
		t.Fatalf("Error creating LogListener: %v", err)
	}

	consoleChannel := make(chan logbus.LogEvent, 10)
	consoleConsumer := &PrintToConsole{Name: "ConsoleConsumer", Events: consoleChannel}

	go consoleConsumer.Consume()

	broadcaster.Subscribe(consoleChannel)

	go listener.StartSimulateListening()

	time.Sleep(5 * time.Second)

	log.Printf("🐳 Done")
}

func TestLogListenerWithPrintToConsole(t *testing.T) {
	broadcaster := logbus.NewLogBroadcaster()

	chConsole := make(chan logbus.LogEvent, 10)
	consoleConsumer := &PrintToConsole{Name: "ConsoleConsumer", Events: chConsole}

	go consoleConsumer.Consume()

	broadcaster.Subscribe(chConsole)

	config := &ListenerConfig{WebSocketURL: "ws://localhost:8546"}
	decoder := &DefaultDecoder{}
	listener, err := NewLogListener(config, broadcaster, decoder)
	if err != nil {
		t.Fatalf("❌ Failed to create LogListener: %v", err)
	}

	go listener.Listen(context.Background())

	time.Sleep(300 * time.Second)

}

func TestLogListenerWithDecoders(t *testing.T) {
	broadcaster := logbus.NewLogBroadcaster()

	chConsole := make(chan logbus.LogEvent, 10)
	consoleConsumer := &PrintToConsole{Name: "ConsoleConsumer", Events: chConsole}

	go consoleConsumer.Consume()

	broadcaster.Subscribe(chConsole)

	config := &ListenerConfig{WebSocketURL: "ws://localhost:8546"}
	decoder := &DefaultDecoder{
		DecoderFn: GetDecoderFn(),
	}

	listener, err := NewLogListener(config, broadcaster, decoder)
	if err != nil {
		t.Fatalf("❌ Failed to create LogListener: %v", err)
	}

	go listener.Listen(context.Background())

	time.Sleep(300 * time.Second)

}
