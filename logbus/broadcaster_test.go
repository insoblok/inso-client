package logbus

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"time"
)
import "testing"

type Consumer struct {
	Name             string        // Name of the consumer (used in logs)
	Delay            time.Duration // Delay to simulate processing time
	ConsumptionCount int           // Counter to track how many events the consumer processed
	Channel          chan LogEvent // Channel to receive events
}

func startConsumer(name string, delay time.Duration, capacity int, done chan struct{}) *Consumer {
	log.Printf("%s: Starting consumer with delay %v: seconds", name, delay.Seconds())
	ch := make(chan LogEvent, capacity) // buffered channel with capacity for 1 event

	// Start a goroutine for the consumer
	go func() {
		for {
			select {
			case event := <-ch:
				log.Printf("✅ %s: Picked up event %v with TxHash: %s", name, event.LogType, event.TxHash)
				time.Sleep(delay) // simulate work (processing delay)
				log.Printf("🧈 %s: Consumed event %v with TxHash: %s", name, event.LogType, event.TxHash)
			case <-done: // Close the goroutine once done channel is closed
				log.Printf("%s: Exiting", name)
				return
			}
		}
	}()

	return &Consumer{
		Name:    name,
		Delay:   delay,
		Channel: ch,
	}
}

// PublishEventsInLoop will publish events in a loop with a delay between each one.
func PublishEventsInLoop(b LogBroadcaster, n int, waitTime time.Duration) {
	for i := 0; i < n; i++ {
		// Create an event with an incremented event number
		event := LogEvent{
			LogType: UnknownEventLog,
			TxHash:  fmt.Sprintf("0x%x", i), // Unique txHash as a hex string
		}

		// Log the publishing of the event
		log.Printf("🎤 Publishing event %d with TxHash: %s", i, event.TxHash)
		b.Publish(event) // Publish to all subscribers

		// Log the event publishing completion
		log.Printf("🐳 Published event %d with TxHash: %s", i, event.TxHash)

		// Wait for `waitTime` seconds before publishing the next event
		time.Sleep(waitTime)
	}
}

/////////////////////////////

func TestBroadcastCreation(t *testing.T) {
	b := NewLogBroadcaster() // ✅ returns interface
	assert.NotNil(t, b)
}

func TestSubscribeAndPublish(t *testing.T) {
	b := NewLogBroadcaster()

	ch := make(chan LogEvent, 1)
	b.Subscribe(ch)

	event := LogEvent{
		LogType: UnknownEventLog,
		TxHash:  "0x123",
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
		LogType: UnknownEventLog,
		TxHash:  "0xabc",
	}

	b.Publish(event)

	assert.Equal(t, event.TxHash, (<-ch1).TxHash)
	assert.Equal(t, event.TxHash, (<-ch2).TxHash)
}

func TestSlowSubscriberDoesNotBlockOthers(t *testing.T) {
	b := NewLogBroadcaster()

	// Create done channel to notify when the consumers are done
	done := make(chan struct{})

	// Create slow and fast consumers
	chSlow := startConsumer("SlowConsumer", 10*time.Second, 1, done)
	chFast := startConsumer("FastConsumer", 1*time.Second, 1, done)

	b.Subscribe(chSlow.Channel)
	b.Subscribe(chFast.Channel)

	log.Printf("🎤 Publishing E1")
	event1 := LogEvent{LogType: UnknownEventLog, TxHash: "0x1"}
	b.Publish(event1) // Goes to both chSlow and chFast
	log.Printf("🐳 Published E1")
	// Allow some time for consumers to process the event
	time.Sleep(2 * time.Second)

	log.Printf("🎤 Publishing E2")
	event2 := LogEvent{LogType: UnknownEventLog, TxHash: "0x2"}
	b.Publish(event2)
	log.Printf("🐳Published E2")

	// Allow chFast to receive event2
	received := <-chFast.Channel // chFast should get event2
	assert.Equal(t, "0x2", received.TxHash)

	// Ensure chSlow is skipped and hasn't received anything
	select {
	case <-chSlow.Channel:
		t.Error("Expected chSlow to be skipped")
	default:
		// chSlow is expected to not receive anything
	}

	close(done)
}

func TestPublisherLoop(t *testing.T) {
	b := NewLogBroadcaster()

	done := make(chan struct{})
	chSlow := startConsumer("SlowConsumer", 5*time.Second, 1, done)
	chFast := startConsumer("FastConsumer", 1*time.Second, 1, done)

	b.Subscribe(chSlow.Channel)
	b.Subscribe(chFast.Channel)

	PublishEventsInLoop(b, 5, 2*time.Second)

	log.Printf("🛌 Publisher go to sleep for 60 seconds")
	time.Sleep(60 * time.Second)
	log.Printf("👋 App Exiting")
	close(done)
}

func TestPublisherLoop2(t *testing.T) {
	b := NewLogBroadcaster()

	done := make(chan struct{})
	chSlow := startConsumer("SlowConsumer", 10*time.Second, 5, done)
	chFast := startConsumer("FastConsumer", 1*time.Second, 1, done)

	b.Subscribe(chSlow.Channel)
	b.Subscribe(chFast.Channel)

	PublishEventsInLoop(b, 5, 2*time.Second)

	log.Printf("🛌 Publisher go to sleep for 60 seconds")
	time.Sleep(60 * time.Second)
	log.Printf("👋 App Exiting")
	close(done)
}
