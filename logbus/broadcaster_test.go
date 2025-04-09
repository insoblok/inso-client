package logbus

import "github.com/stretchr/testify/assert"
import "testing"

func TestBroadcastCreation(t *testing.T) {
	b := NewLogBroadcaster() // ✅ returns interface
	assert.NotNil(t, b)
}

func TestSubscribeAndPublish(t *testing.T) {
	b := NewLogBroadcaster()

	ch := make(chan LogEvent, 1)
	b.Subscribe(ch)

	event := LogEvent{
		Event:  "TestEvent",
		TxHash: "0x123",
	}

	b.Publish(event)

	received := <-ch
	assert.Equal(t, event.TxHash, received.TxHash)
}
