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

func TestDuplicateSubscriptionPanics(t *testing.T) {
	b := NewLogBroadcaster()
	ch := make(chan LogEvent, 1)
	b.Subscribe(ch)
	assert.Panics(t, func() {
		b.Subscribe(ch)
	})
}

func TestMultipleSubscribers(t *testing.T) {
	b := NewLogBroadcaster()

	ch1 := make(chan LogEvent, 1)
	ch2 := make(chan LogEvent, 1)

	b.Subscribe(ch1)
	b.Subscribe(ch2)

	event := LogEvent{
		Event:  "TestEventMulti",
		TxHash: "0xabc",
	}

	b.Publish(event)

	assert.Equal(t, event.TxHash, (<-ch1).TxHash)
	assert.Equal(t, event.TxHash, (<-ch2).TxHash)
}

func TestSlowSubscriberDoesNotBlockOthers(t *testing.T) {
	b := NewLogBroadcaster()

	chSlow := make(chan LogEvent, 1)
	chFast := make(chan LogEvent, 2)

	b.Subscribe(chSlow)
	b.Subscribe(chFast)

	event1 := LogEvent{Event: "E1", TxHash: "0x1"}
	event2 := LogEvent{Event: "E2", TxHash: "0x2"}

	// Publish event 1
	b.Publish(event1) // goes to both
	// Do not read from chSlow (simulate slow subscriber)

	// Publish event 2
	b.Publish(event2) // chSlow is full now, should skip, but chFast should get it

	// Allow chFast to receive event2
	received := <-chFast // chFast should get event2
	assert.Equal(t, "0x2", received.TxHash)

	// Ensure chSlow is skipped and hasn't received anything
	select {
	case <-chSlow:
		t.Error("Expected chSlow to be skipped")
	default:
		// chSlow is expected to not receive anything
	}
}

func TestSlowSubscriberDoesNotBlockOthers2(t *testing.T) {
	b := NewLogBroadcaster()

	chSlow := make(chan LogEvent, 1) // Slow subscriber with a buffer of 1
	chFast := make(chan LogEvent, 2) // Fast subscriber with a buffer of 2

	b.Subscribe(chSlow)
	b.Subscribe(chFast)

	// Event 1 - Publish to both channels
	event1 := LogEvent{Event: "E1", TxHash: "0x1"}
	b.Publish(event1) // Goes to both chSlow and chFast

	// Do not read from chSlow (simulate slow subscriber)
	// chSlow now has event1 and is full

	// Event 2 - Publish to both channels
	event2 := LogEvent{Event: "E2", TxHash: "0x2"}
	b.Publish(event2) // chSlow is full now, should skip, but chFast should get it

	// Ensure that chFast receives event2
	received := <-chFast
	assert.Equal(t, "0x2", received.TxHash)

	// Ensure chSlow is skipped and hasn't received anything
	select {
	case <-chSlow:
		t.Error("Expected chSlow to be skipped")
	default:
		// chSlow should not have received any events
	}
}
