package bus

import (
	"sync"

	"github.com/google/uuid"
)

// Bus provides message pub/sub functionality via Publish() method
type Bus struct {
	mu               sync.RWMutex
	subscribers      map[string]*Subscriber
	subscriberBuffer int
}

// Subscriber belongs to Bus
type Subscriber struct {
	id string
	ch chan any
}

func New(subscriberBuffer int) *Bus {
	return &Bus{
		mu:               sync.RWMutex{},
		subscribers:      make(map[string]*Subscriber),
		subscriberBuffer: subscriberBuffer,
	}
}

func (b *Bus) Subscribe() *Subscriber {
	sub := &Subscriber{
		id: uuid.New().String(),
		ch: make(chan any, b.subscriberBuffer),
	}

	b.mu.Lock()
	b.subscribers[sub.id] = sub
	b.mu.Unlock()

	return sub
}

func (b *Bus) Unsubscribe(sub *Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// safety check to prevent closing of closed channel
	if _, exists := b.subscribers[sub.id]; exists {
		close(sub.ch)
		delete(b.subscribers, sub.id)
	}
}

func (b *Bus) Publish(message any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.subscribers {
		select {
		case sub.ch <- message:
		default: // safe skip if subscriber buffer is full
		}
	}
}

func (b *Bus) Subscribers() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.subscribers)
}

func (s *Subscriber) Channel() <-chan any {
	return s.ch
}

func (s *Subscriber) ID() string {
	return s.id
}
