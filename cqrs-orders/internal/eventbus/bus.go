package eventbus

import (
	"cqrs-orders/internal/domain"
	"sync"
)

// Handler is any function that can process a domain event
type Handler func(domain.Event)

// Bus is an in-memory pub/sub channel.
// Publishers don't know who's listening - subscribers don't know who's publishing
type Bus struct {
	mu          sync.RWMutex
	subscribers []Handler
	ch          chan domain.Event
}

func New(bufferSize int) *Bus {
	return &Bus{
		ch:          make(chan domain.Event, bufferSize),
		subscribers: make([]Handler, 0),
	}
}

// Subscribe registers a handler that will be called for every event.
// Call this before Start - subscribers registered after Start may miss early events
func (b *Bus) Subscribe(h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, h)
}

// Publish sends an event to the channel. Non-blocking up to buffer capacity.
// If the channel is full, this will blosk - in production you'd handle this
// with a timeout or dead-letter queue.
func (b *Bus) Publish(e domain.Event) {
	b.ch <- e
}

// Start launches the dispatch loop in a background goroutine.
// It reads events from the channel and fans them out to all subscribers.
// Call Stop when shutting down to drain the channel cleanly
func (b *Bus) Start() {
	go func() {
		for e := range b.ch {
			b.mu.RLock()
			subs := make([]Handler, len(b.subscribers))
			copy(subs, b.subscribers)
			b.mu.RUnlock()

			// Each subscriber gets the event
			// In production, you'd run these concurrently with individual error handling
			for _, sub := range subs {
				sub(e)
			}
		}
	}()
}

// Stop closes the channel, which causes the dispatch loop to exit
// after processing all buffered events. Always call this on shutdown
// or your last events will be lost
func (b *Bus) Stop() {
	close(b.ch)
}
