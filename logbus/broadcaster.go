package logbus

import "sync"

type LogEvent struct {
	Contract  string                 // Optional: contract name or address
	TxHash    string                 // Transaction hash
	Block     uint64                 // Block number
	Timestamp int64                  // Optional: unix time
	Args      map[string]interface{} // Decoded args (message, value, etc.)
	LogType   LogType                // Log type (Transaction, Event, etc.)
}

type LogBroadcaster interface {
	Subscribe(chan<- LogEvent)
	Unsubscribe(chan<- LogEvent)
	Publish(LogEvent)
}

type inMemoryBroadcaster struct {
	subscribers []chan<- LogEvent
	mu          sync.Mutex
}

func NewLogBroadcaster() LogBroadcaster {
	return &inMemoryBroadcaster{
		subscribers: make([]chan<- LogEvent, 0),
	}
}

func (b *inMemoryBroadcaster) Subscribe(ch chan<- LogEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subscribers {
		if sub == ch {
			panic("channel already subscribed") // or return an error/log
		}
	}

	b.subscribers = append(b.subscribers, ch)
}

func (b *inMemoryBroadcaster) Publish(event LogEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subscribers {
		select {
		case ch <- event:
			// sent successfully
		default:
			// optional: drop or log if subscriber is full
		}
	}
}

func (b *inMemoryBroadcaster) Unsubscribe(ch chan<- LogEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, sub := range b.subscribers {
		if sub == ch {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			break
		}
	}
}
